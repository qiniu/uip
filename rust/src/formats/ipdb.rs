use std::io;
use crate::formats::{Query, QueryBuilder};
use crate::data::{VersionInfo, IpInfo};
mod reader;

struct IpdbBuilder {
}

impl QueryBuilder for IpdbBuilder {
    fn build(bytes: &[u8]) -> io::Result<Box<dyn Query>> {
        let reader = Ipdb::new(bytes)?;
        Ok(Box::new(reader))
    }
}

pub struct Ipdb {
    reader: reader::Reader,
}

impl Ipdb {
    pub fn new(data: &[u8]) -> io::Result<Ipdb> {
        let reader = reader::Reader::new(&data)?;
        Ok(Ipdb { reader })
    }
}

impl Query for Ipdb {
    fn query(&self, ip:&[u8]) -> io::Result<IpInfo>{
        self.reader.find(ip, "CN", true)
    }

    fn build_cache(&mut self, _ip_list:&Vec<String>) {
        println!("ipdb build_cache");
    }

    fn version(&self) -> VersionInfo {
        println!("ipdb get_version");
        VersionInfo {
            ip_type: self.reader.meta.ip_version,
            count: self.reader.meta.node_count,
            build: self.reader.meta.build,
            version: match &self.reader.meta.version {
                None => format!("ipip.net-{}", self.reader.meta.build),
                Some(s) => s.to_string(),
            },
            languages: vec!["CN".into()],
            extra_info: self.reader.meta.extra.clone(),
        }
    }
}