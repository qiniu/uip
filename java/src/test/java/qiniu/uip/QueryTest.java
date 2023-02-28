package qiniu.uip;

import org.junit.Test;
import qiniu.uip.db.IPFormatException;
import qiniu.uip.db.IpInfo;
import qiniu.uip.db.Query;

import java.io.IOException;

public class QueryTest {

    @Test
    public void testQueryV4() throws IOException, IPFormatException {
        String file = System.getenv("IPDBv4_PATH");
        file = "/Users/long/github/qiniu/uip/ipv4.ipdb";
        if (file == null) {
            return;
        }
        Query q = Query.create(file);
        IpInfo info = q.query("220.248.53.1");
        System.out.println(info);
        assert info.country.equals("中国");

        assert q.version() != null;
        assert q.toString() != null;
        assert q.version().hasIpV4();
    }

    @Test
    public void testQueryV6() throws IOException, IPFormatException {
        String file = System.getenv("IPDBv6_PATH");
        if (file == null) {
            return;
        }
        Query q = Query.create(file);
        IpInfo info = q.query("2001:4860:4860::8888");
        System.out.println(info);
        assert info.country.equals("美国");

        assert q.version() != null;
        assert q.toString() != null;
        assert q.version().hasIpV6();
    }
}
