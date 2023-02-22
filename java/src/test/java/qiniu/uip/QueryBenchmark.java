package qiniu.uip;

import qiniu.uip.db.IPFormatException;
import qiniu.uip.db.IpInfo;
import qiniu.uip.db.Query;
import qiniu.uip.util.IPAddress;

import java.io.IOException;

public class QueryBenchmark {
    private QueryBenchmark() {
    }

    public static void main(String[] args) throws IOException, IPFormatException {
        String file = args[1];
        if (file == null) {
            return;
        }
        Query q = Query.create(file);
        long start = System.currentTimeMillis();
        int cycle = 1000000;
        byte[] ip = IPAddress.textToNumericFormatV4("220.248.53.61");
        for (int i = 0; i < cycle; i++) {
            IpInfo info = q.query(ip);
        }
        long end = System.currentTimeMillis();
        System.out.println("every query " + (end - start) * 1000000 / cycle + " ns");
    }
}