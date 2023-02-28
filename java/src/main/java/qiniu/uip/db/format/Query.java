package qiniu.uip.db.format;

import qiniu.uip.db.IPFormatException;
import qiniu.uip.db.InvalidDatabaseException;
import qiniu.uip.db.IpInfo;
import qiniu.uip.db.VersionInfo;

public interface Query {
    IpInfo query(byte[] ip) throws IPFormatException, InvalidDatabaseException;

    void buildCache(byte[][] ipList);

    VersionInfo version();
}
