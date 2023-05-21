pub(crate) mod ipdb;

use std::io;
use super::data::VersionInfo;

pub trait Query {
    fn query(&self, ip: &[u8]) -> io::Result<super::data::IpInfo>;
    fn build_cache(&mut self, ip_list:&Vec<String>);
    fn version(&self) -> VersionInfo;
}

pub trait QueryBuilder {
    fn build(bytes: &[u8]) -> io::Result<Box<dyn Query>>;
}

