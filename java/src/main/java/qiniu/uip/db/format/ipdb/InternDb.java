package qiniu.uip.db.format.ipdb;

import qiniu.uip.db.IPFormatException;
import qiniu.uip.db.InvalidDatabaseException;
import qiniu.uip.db.IpInfo;
import qiniu.uip.db.VersionInfo;
import qiniu.uip.db.format.Query;
import qiniu.uip.db.format.QueryBuilder;
import qiniu.uip.util.Strings;

public final class InternDb implements Query {

    private final Reader reader;

    public InternDb(byte[] dat) throws InvalidDatabaseException {
        reader = new Reader(dat);
    }

    public IpInfo query(byte[] ip) throws IPFormatException, InvalidDatabaseException {
        String[] parts = reader.find(ip, "CN");
        if (parts == null) {
            throw new IPFormatException("invalid ip address");
        }
        String[] fields = reader.getSupportFields();
        IpInfo info = new IpInfo();
        for (int i = 0, l = parts.length; i < l; i++) {
            switch (fields[i]) {
                case "country_name":
                    info.country = parts[i];
                    break;
                case "region_name":
                    info.province = parts[i];
                    break;
                case "city_name":
                    info.city = parts[i];
                    break;
                case "isp_domain":
                    info.isp = parts[i];
                    break;
                case "asn":
                    info.asn = parts[i];
                case "line":
                    info.line = parts[i];
                    break;
                case "district":
                    info.district = parts[i];
                    break;
                case "continent_code":
                    switch (parts[i]) {
                        case "AS":
                            info.continent = "亚洲";
                            break;
                        case "EU":
                            info.continent = "欧洲";
                            break;
                        case "NA":
                            info.continent = "北美洲";
                            break;
                        case "SA":
                            info.continent = "南美洲";
                            break;
                        case "AF":
                            info.continent = "非洲";
                            break;
                        case "OC":
                            info.continent = "大洋洲";
                            break;
                        case "AN":
                            info.continent = "南极洲";
                            break;
                        default:
                            info.continent = parts[i];
                    }
                    break;
            }
        }
        return info;
    }

    public VersionInfo version() {
        VersionInfo versionInfo = new VersionInfo();
        String version;
        if (Strings.isEmpty(reader.meta.version)) {
            version = "ipip.net-" + reader.meta.build;
        } else {
            version = reader.meta.version;
        }
        versionInfo.ipType = reader.meta.ip_version;
        versionInfo.count = reader.meta.node_count;
        versionInfo.build = reader.meta.build;
        versionInfo.version = version;
        versionInfo.languages = new String[]{"CN"};
        versionInfo.extraInfo = reader.meta.extra;
        return versionInfo;
    }

    public static class Builder implements QueryBuilder {
        public Builder() {
        }

        public InternDb build(byte[] dat) throws InvalidDatabaseException {
            return new InternDb(dat);
        }
    }
}

