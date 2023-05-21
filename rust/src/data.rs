use std::fmt::{Display, Formatter};
use serde::{Deserialize};

pub const IPV4: i32 = 1;
pub const IPV6: i32 = 2;

#[derive(Deserialize)]
// #[derive(Copy, Clone, Display)]
pub struct IpInfo {
    pub country: String,
    pub district: String,
    pub province: String,
    pub city: String,
    pub asn: String,
    pub isp: String,
    pub continent: String,
    pub line: String,
}

impl Display for IpInfo {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "Country: {}, District: {}, Province: {}, City: {}, Asn: {}, Isp: {}, Continent: {}, Line: {}",
            self.country, self.district, self.province, self.city, self.asn, self.isp, self.continent, self.line
        )
    }
}

// impl Copy for IpInfo {
//
// }

pub struct VersionInfo {
    pub ip_type: i32,
    pub count: i32,
    pub build: i64,
    pub version: String,
    pub languages: Vec<String>,
    pub extra_info: Option<Vec<String>>,
}

impl VersionInfo {
    pub fn has_ipv4(&self) -> bool {
        self.ip_type & IPV4 != 0
    }

    pub fn has_ipv6(&self) -> bool {
        self.ip_type & IPV6 != 0
    }
}

impl Display for VersionInfo {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "ip_type: {}, count: {}, build: {}, version: {}, languages: {:?}, extra_info: {:?}",
            self.ip_type, self.count, self.build, self.version, self.languages, self.extra_info
        )
    }
}