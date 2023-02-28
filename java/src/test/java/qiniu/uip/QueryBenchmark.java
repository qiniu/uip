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
        String file = System.getenv("IPDBv4_PATH");
        if (file == null) {
            return;
        }
        Query q = Query.create(file);
        int cycle = 1000000;
        byte[] ip = IPAddress.textToNumericFormatV4("220.248.53.1");
        IpInfo info = q.query(ip);
        System.out.println(info);
        //warm up
        for (int i = 0; i < 100; i++) {
            info = q.query(ip);
        }

        long start = System.currentTimeMillis();

        for (int i = 0; i < cycle; i++) {
            info = q.query(ip);
        }
        long end = System.currentTimeMillis();
        System.out.println("every query " + (end - start) * 1000000 / cycle + " ns");
    }
}
