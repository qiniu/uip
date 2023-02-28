package qiniu.uip.db.format.ipdb;

import qiniu.uip.db.IPFormatException;
import qiniu.uip.db.InvalidDatabaseException;
import qiniu.uip.db.IpInfo;
import qiniu.uip.db.VersionInfo;
import qiniu.uip.db.format.Query;
import qiniu.uip.db.format.QueryBuilder;
import qiniu.uip.util.Strings;

public final class Ipdb implements Query {

    private final Reader reader;

    public Ipdb(byte[] dat) throws InvalidDatabaseException {
        reader = new Reader(dat);
    }


    @Override
    public IpInfo query(byte[] ip) throws IPFormatException, InvalidDatabaseException {
        return reader.query(ip, true, false);
    }

    @Override
    public IpInfo query(long ab, long cd) throws IPFormatException, InvalidDatabaseException {
        return null;
    }

    @Override
    public void buildCache(byte[][] ipList) {
        reader.buildCache(ipList);
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

        public Ipdb build(byte[] dat) throws InvalidDatabaseException {
            return new Ipdb(dat);
        }
    }
}

