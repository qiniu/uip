package qiniu.uip.db.format;

import qiniu.uip.db.IPFormatException;
import qiniu.uip.db.InvalidDatabaseException;
import qiniu.uip.db.IpInfo;
import qiniu.uip.db.VersionInfo;

public interface Query {
    IpInfo query(byte[] ip) throws IPFormatException, InvalidDatabaseException;

    // ipv6 128bit = 2 long, avoid object allocation
    IpInfo query(long ab, long cd) throws IPFormatException, InvalidDatabaseException;

    void buildCache(byte[][] ipList);

    VersionInfo version();
}
