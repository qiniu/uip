package qiniu.uip.db;

public final class IpInfo {
    public String country;
    public String district;
    public String province;
    public String city;
    public String asn;
    public String isp;
    public String continent;
    public String line;

    public String toString() {
        return String.format("country: %s, district: %s, province: %s, city: %s, asn: %s, isp: %s, continent: %s, line: %s", country, district, province, city, asn, isp, continent, line);
    }
}
