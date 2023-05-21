use std::collections::HashMap;
use std::io;
use std::io::{Result, Error, ErrorKind};

use byteorder::{ByteOrder, BigEndian};

use serde::{Deserialize};

use crate::data::{IpInfo};

#[derive(Deserialize)]
pub struct MetaData {
    #[serde(skip_deserializing)]
    pub meta_size: usize,
    pub build: i64,
    pub ip_version: i32,
    pub node_count: i32,
    pub languages: HashMap<String, i32>,
    pub fields: Vec<String>,
    pub total_size: i32,
    pub version: Option<String>,
    pub extra: Option<Vec<String>>,
}

impl MetaData {
    fn new(buf: &[u8]) -> Result<MetaData> {
        if buf.len() < 4 {
            return Err(Error::new(ErrorKind::InvalidData, "meta data length < 4"));
        }

        let meta_size = BigEndian::read_i32(&buf[0..4]) as usize;
        let end = meta_size + 4;
        if buf.len() < end {
            return Err(Error::new(ErrorKind::InvalidData, "meta data length < meta_size"));
        }
        let v:serde_json::Result<MetaData> = serde_json::from_slice(&buf[4..end]);
        return match v {
            Ok(mut m) => {
                m.meta_size = meta_size;
                Ok(m)
            },
            Err(e) => {
                Err(Error::new(ErrorKind::InvalidData, e))
            }
        }
    }
}

const CACHE_DEPTH: i32 = 12;

pub struct Reader {
    pub meta: MetaData,
    file_size: i32,
    node_count: i32,
    data: Vec<u8>,
    v4_offset: i32,
    ipv4_bits_cache: Vec<i32>,
    _result_cache: HashMap<i32, IpInfo>,
}

impl Reader{
    pub fn new(bytes: &[u8]) -> Result<Reader> {
        let m = MetaData::new(&bytes)?;
        if m.total_size + m.meta_size as i32 +4 != bytes.len() as i32 {
            return Err(Error::new(ErrorKind::InvalidData, "invalid data length"));
        }
        let data = (&bytes[m.meta_size + 4..]).to_vec();

        let mut reader = Reader {
            file_size: bytes.len() as i32,
            node_count: m.node_count,
            v4_offset: 0,
            ipv4_bits_cache: Vec::with_capacity(1<<CACHE_DEPTH),
            _result_cache: HashMap::new(),
            data,
            meta: m,
        };
        if reader.is_ipv4() {
            let mut node = 0;
            for i in 0..96 {
                if i >= 80 {
                    node = reader.read_node(node, 1);
                } else {
                    node = reader.read_node(node, 0);
                }
                if node >= reader.node_count {
                    break;
                }
            }
            reader.v4_offset = node;
            reader.init_cache()
        } else {
            reader.v4_offset = 0
        }
        Ok(reader)
    }

    fn read_node(&self, node: i32, index:i32) -> i32 {
        let off = node * 8 + index * 4;
        BigEndian::read_i32(&self.data[off as usize..off as usize + 4])
    }

    fn init_cache(&mut self) {
        //construct cache from binary trie tree for reduce read memory time
        let mut b: [u8;2] = [0,0];
        for i in 0..(1<<CACHE_DEPTH) {
            // index to bytes
            b[0] = ((i << (16 - CACHE_DEPTH)) >> 8) as u8;
            b[1] = (0xFF & (i << (16 - CACHE_DEPTH))) as u8;

            let node = self.read_depth(self.v4_offset, CACHE_DEPTH, 0, &b);
            self.ipv4_bits_cache.push(node);
        }
    }

    fn is_ipv4(&self) -> bool {
        self.meta.ip_version & 1 == 1
    }

    fn is_ipv6(&self) -> bool {
        self.meta.ip_version & 2 == 2
    }

    pub fn find(&self, addr :&[u8], lang:&str, _cache:bool) -> Result<IpInfo> {
        let node = self.find0(addr)?;
        // if cache && self.result_cache.len() > 0 {
        //     let x = self.result_cache.get(&node);
        //     match x {
        //         Some(i) => {
        //             return Ok(*i);
        //         }
        //         None => {}
        //     }
        // }
        let v = self.resolve_node(node, lang)?;
        Ok(Self::build_info(&v, &self.meta.fields))
    }

    fn build_info(val:&Vec<String>, fields: &Vec<String>) -> IpInfo {
        let mut info = IpInfo{
            country: "".to_string(),
            district: "".to_string(),
            province: "".to_string(),
            city: "".to_string(),
            asn: "".to_string(),
            isp: "".to_string(),
            continent: "".to_string(),
            line: "".to_string(),
        };
        for i in 0..fields.len() {
            match fields[i].as_str(){
                "country_name"       => info.country = val[i].to_string(),
                "region_name" => info.province = val[i].to_string(),
                "city_name" => info.city = val[i].to_string(),
                "isp_domain"=> info.isp = val[i].to_string(),
                "asn" => info.asn = val[i].to_string(),
                "line" => info.line = val[i].to_string(),
                "district" => info.district = val[i].to_string(),
                "continent_code" => {
                    match val[i].as_str() {
                        "AS" => info.continent = "亚洲".to_string(),
                        "EU" => info.continent = "欧洲".to_string(),
                        "NA" => info.continent = "北美洲".to_string(),
                        "SA" => info.continent = "南美洲".to_string(),
                        "AF" => info.continent = "非洲".to_string(),
                        "OC" => info.continent = "大洋洲".to_string(),
                        "AN" => info.continent = "南极洲".to_string(),
                        _ => {}
                    }
                },
                _ => {}
            };
        }
        info
    }

    fn find0(&self, ip: &[u8]) -> Result<i32> {
        if ip.len() == 16 {
            if !self.is_ipv6() {
                return Err(Error::new(ErrorKind::InvalidData, "no support ipv6"));
            }
        } else if ip.len() == 4 {
            if !self.is_ipv4() {
                return Err(Error::new(ErrorKind::InvalidData, "no support ipv4"));
            }
        } else {
            return Err(Error::new(ErrorKind::InvalidData, "invalid ip"));
        }
        self.find_node(ip)
    }

    fn resolve_node(&self, node: i32, lang:&str) -> Result<Vec<String>> {
        let off = self.meta.languages.get(lang);
        match off {
            None => return Err(Error::new(ErrorKind::InvalidData, "invalid language")),
            Some(_) => {
                let resolved = node - self.node_count + self.node_count * 8;
                if resolved >= self.file_size {
                    return Err(Error::new(ErrorKind::InvalidData, "database resolve error"));
                }

                let resolved: usize = resolved as usize;

                let size = BigEndian::read_u16(&self.data[resolved ..resolved + 2]) as usize;

                if self.data.len() < resolved + 2 + size {
                    return Err(Error::new(ErrorKind::InvalidData, "database resolve error"));
                }
                let str = String::from_utf8_lossy(&self.data[resolved  + 2..resolved  + 2 + size ]);
                let a = str.split("\t");
                let r = a.map(|i| i.to_string()).collect();
                Ok(r)
            }
        }
    }

    fn find_node(&self, b: &[u8]) -> Result<i32> {
        let mut node = 0;
        let bit = b.len() * 8;
        if bit == 32 {
            node = self.v4_offset;
        }
        let mut i = 0;
        if self.ipv4_bits_cache.len()>0 && bit == 32 {
            let index = Self::bytes_to_index(b);
            node = self.ipv4_bits_cache[index as usize];
            if node >= self.node_count {
                return Ok(node);
            }
            i = CACHE_DEPTH;
        }
        node = self.read_depth(node, bit as i32, i, b);
        if node >= self.node_count {
            return Ok(node);
        }

        return Err(io::ErrorKind::NotFound.into());
    }

    fn bytes_to_index(b: &[u8]) -> i32 {
        let i1 = ((0xFF & (b[0] as i32)) << 8) >> (16 - CACHE_DEPTH);
        let i2 = (0xFF & (b[1] as i32)) >> (16 - CACHE_DEPTH);
        i1 | i2
    }

    fn read_depth(&self, node0: i32, depth:i32, start:i32, b:&[u8]) -> i32 {
        let mut node = node0;
        for i in start..depth {
            if node >= self.node_count {
                break;
            }
            let index = 1 & ((0xFF & b[(i / 8) as usize] as usize) >> 7 - (i % 8));
            node = self.read_node(node, index as i32);
        }
        node
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn meta_deser() {

    }
}
