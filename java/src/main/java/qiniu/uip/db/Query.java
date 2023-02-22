package qiniu.uip.db;

import qiniu.uip.db.format.Formats;
import qiniu.uip.db.format.QueryBuilder;
import qiniu.uip.db.format.ipdb.InternDb;
import qiniu.uip.util.File;
import qiniu.uip.util.IPAddress;
import qiniu.uip.util.Strings;

import java.io.IOException;

public final class Query {
    static {
        Formats.registerQueryFormat(".ipdb", new InternDb.Builder());
    }

    private qiniu.uip.db.format.Query q;

    private Query(qiniu.uip.db.format.Query q) {
        this.q = q;
    }

    public static Query create(String filename) throws IOException {
        String ext = filename.substring(filename.lastIndexOf("."));
        byte[] data = File.readAll(filename);
        return create(data, ext);
    }

    public static Query create(byte[] data, String format) throws InvalidDatabaseException {
        QueryBuilder qb = Formats.getQueryFormat(format);
        if (qb == null) {
            throw new InvalidDatabaseException("not support format");
        }
        return new Query(qb.build(data));
    }

    public IpInfo query(String addr) throws IPFormatException, InvalidDatabaseException {
        byte[] ipv;
        if (Strings.isEmpty(addr)) {
            throw new IPFormatException("ip is empty");
        }

        if (addr.indexOf(":") > 0) {
            ipv = IPAddress.textToNumericFormatV6(addr);
            if (ipv == null) {
                throw new IPFormatException("ipv6 format error");
            }

        } else if (addr.indexOf(".") > 0) {
            ipv = IPAddress.textToNumericFormatV4(addr);
            if (ipv == null) {
                throw new IPFormatException("ipv4 format error");
            }
        } else {
            throw new IPFormatException("ip format error");
        }
        return query(ipv);
    }

    public IpInfo query(byte[] addr) throws IPFormatException, InvalidDatabaseException {
        return q.query(addr);
    }

    public VersionInfo version() {
        return q.version();
    }
}
