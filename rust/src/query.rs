use std::io::Read;
use super::formats::Query;
use super::formats::ipdb::Ipdb;
use std::net::{IpAddr};

pub struct QueryDb {
    query: Box<dyn Query>,
}

impl QueryDb {
    pub fn from_file(path: &str) -> std::io::Result<QueryDb> {
        let mut file = std::fs::File::open(path)?;
        let mut buffer = Vec::new();
        file.read_to_end(&mut buffer)?;
        Self::from_bytes(&buffer)
    }

    pub fn from_bytes(bytes: &[u8]) -> std::io::Result<QueryDb> {
        let query = Ipdb::new(bytes)?;
        let b = Box::new(query);
        Ok(QueryDb { query: b })
    }

    pub fn query(&self, ip: &[u8]) -> std::io::Result<super::data::IpInfo> {
        self.query.query(ip)
    }

    pub fn query_str(&self, ip: &str) -> std::io::Result<super::data::IpInfo> {
        return match ip.parse::<IpAddr>() {
            Ok(ipa) => {
                match ipa {
                    IpAddr::V4(ipv4) => {
                        self.query.query(&ipv4.octets())
                    },
                    IpAddr::V6(ipv6) => {
                        self.query.query(&ipv6.octets())
                    },
                }
            },
            Err(_) => {
                Err(std::io::Error::new(std::io::ErrorKind::Other, "parse ip error"))
            },
        }
    }

    pub fn query_u32(&self, ip: u32) -> std::io::Result<super::data::IpInfo> {
        let ipb = ip.to_be_bytes();
        self.query(&ipb)
    }

    // pub fn build_cache(&mut self, ip_list: &Vec<String>) {
    //     self.query.build_cache(ip_list);
    // }

    pub fn version(&self) -> super::data::VersionInfo {
        self.query.version()
    }
}

