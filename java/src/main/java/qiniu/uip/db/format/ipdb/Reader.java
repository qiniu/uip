package qiniu.uip.db.format.ipdb;

import qiniu.uip.db.IPFormatException;
import qiniu.uip.db.InvalidDatabaseException;
import qiniu.uip.db.IpInfo;
import qiniu.uip.db.NotFoundException;
import qiniu.uip.util.Json;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.HashMap;


final class Reader {

    private static final int cacheDepth = 12;
    MetaData meta;
    private int fileSize;
    private int nodeCount;
    private byte[] data;
    private int v4offset;
    private int[] ipv4BitsCache;

    private HashMap<Integer, IpInfo> resultCache;

    Reader(byte[] data) throws InvalidDatabaseException {
        this.data = data;
        this.fileSize = data.length;
        if (this.fileSize < 5) {
            throw new InvalidDatabaseException("database file size error");
        }

        long metaLength = bytesToLong(
                this.data[0],
                this.data[1],
                this.data[2],
                this.data[3]
        );

        try {
            int end = Long.valueOf(metaLength).intValue() + 4;
            byte[] metaBytes = Arrays.copyOfRange(this.data, 4, end);

            MetaData meta = Json.parseObject(new String(metaBytes), MetaData.class);

            this.nodeCount = meta.node_count;
            this.meta = meta;
        } catch (Exception e) {
            throw new InvalidDatabaseException(e.getMessage());
        }

        if ((meta.total_size + Long.valueOf(metaLength).intValue() + 4) != this.data.length) {
            throw new InvalidDatabaseException("database file size error");
        }

        this.data = Arrays.copyOfRange(this.data, Long.valueOf(metaLength).intValue() + 4, this.fileSize);

        /** for ipv4 */
        if (0x01 == (this.meta.ip_version & 0x01)) {
            int node = 0;
            for (int i = 0; i < 96 && node < this.nodeCount; i++) {
                if (i >= 80) {
                    node = this.readNode(node, 1);
                } else {
                    node = this.readNode(node, 0);
                }
            }
            this.v4offset = node;
            initCache();
        } else {
            this.ipv4BitsCache = null;
            this.v4offset = 0;
        }
    }

    private static byte[] indexToBytes(int i) {
        int i1 = (i << (16 - cacheDepth)) >> 8;
        int i2 = 0xFF & (i << (16 - cacheDepth));
        return new byte[]{(byte) i1, (byte) i2};
    }

    private static long bytesToLong(byte a, byte b, byte c, byte d) {
        return int2long((((a & 0xff) << 24) | ((b & 0xff) << 16) | ((c & 0xff) << 8) | (d & 0xff)));
    }

    private static long int2long(int i) {
        long l = i & 0x7fffffffL;
        if (i < 0) {
            l |= 0x080000000L;
        }
        return l;
    }

    void buildCache(byte[][] ipList) {
        this.resultCache = new HashMap<>();
        for (byte[] ip : ipList) {
            try {
                query(ip, false, true);
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
    }

    IpInfo query(byte[] ip, boolean cache, boolean build) throws IPFormatException, InvalidDatabaseException {
        int node = find0(ip);
        if (cache && this.resultCache != null) {
            IpInfo info = this.resultCache.get(node);
            if (info != null) {
                return info;
            }
        }
        String[] parts = resolveNode(node, "CN");
        if (parts == null) {
            throw new IPFormatException("invalid ip address");
        }
        String[] fields = getSupportFields();
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
        if (build && this.resultCache != null) {
            this.resultCache.put(node, info);
        }
        return info;
    }

    private int bytesToIndex(byte[] b) {
        int i1 = (0xFF & (int) (b[0])) << 8 >> (16 - cacheDepth);
        int i2 = (0xFF & (int) (b[1])) >> (16 - cacheDepth);
        return i1 | i2;
    }

    private void initCache() {
        ipv4BitsCache = new int[1 << cacheDepth];
        //construct cache from binary trie tree for reduce read memory time
        for (int i = 0; i < ipv4BitsCache.length; i++) {
            byte[] b = indexToBytes(i);
            int node = readDepth(v4offset, cacheDepth, 0, b);
            ipv4BitsCache[i] = node;
        }
    }

    int find0(byte[] ip) throws IPFormatException, InvalidDatabaseException {
        if (ip.length == 16) {
            if (!isIPv6()) {
                throw new IPFormatException("no support ipv6");
            }
        } else if (ip.length == 4) {
            if (!isIPv4()) {
                throw new IPFormatException("no support ipv4");
            }
        } else {
            throw new IPFormatException("ip format error");
        }

        // find node
        int node = 0;
        try {
            node = this.findNode(ip);
        } catch (NotFoundException nfe) {
            return 0;
        }

        return node;
    }

    String[] find(byte[] addr, String language) throws IPFormatException, InvalidDatabaseException {
        int node = find0(addr);

        return resolveNode(node, language);
    }

    private String[] resolveNode(int node, String language) throws InvalidDatabaseException {
        int off;
        try {
            off = this.meta.languages.get(language);
        } catch (NullPointerException e) {
            return null;
        }

        final int resolved = node - this.nodeCount + this.nodeCount * 8;
        if (resolved >= this.fileSize) {
            throw new InvalidDatabaseException("database resolve error");
        }

        byte b = 0;
        int size = Long.valueOf(bytesToLong(
                b,
                b,
                this.data[resolved],
                this.data[resolved + 1]
        )).intValue();

        if (this.data.length < (resolved + 2 + size)) {
            throw new InvalidDatabaseException("database resolve error");
        }

        final String data = new String(this.data, resolved + 2, size, StandardCharsets.UTF_8);
        return Arrays.copyOfRange(data.split("\t", this.meta.fields.length * this.meta.languages.size()), off, off + this.meta.fields.length);
    }

    private int readDepth(int node, int depth, int i, byte[] binary) {
        for (; i < depth; i++) {
            if (node >= this.nodeCount) {
                break;
            }
            node = this.readNode(node, 1 & ((0xFF & binary[i / 8]) >> 7 - (i % 8)));
        }
        return node;
    }

    private int findNode(byte[] binary) throws NotFoundException {

        int node = 0;

        final int bit = binary.length * 8;

        if (bit == 32) {
            node = this.v4offset;
        }
        int i = 0;
        if (ipv4BitsCache != null && bit == 32) {
            int index = bytesToIndex(binary);
            node = ipv4BitsCache[index];
            if (node > this.nodeCount) {
                return node;
            }
            i = cacheDepth;
        }
        node = readDepth(node, bit, i, binary);
        if (node > this.nodeCount) {
            return node;
        }

        throw new NotFoundException("ip not found");
    }

    private int readNode(int node, int index) {
        int off = node * 8 + index * 4;

        return Long.valueOf(bytesToLong(
                this.data[off],
                this.data[off + 1],
                this.data[off + 2],
                this.data[off + 3]
        )).intValue();
    }

    boolean isIPv4() {
        return (this.meta.ip_version & 0x01) == 0x01;
    }

    boolean isIPv6() {
        return (this.meta.ip_version & 0x02) == 0x02;
    }

    int getBuildUTCTime() {
        return this.meta.build;
    }

    String[] getSupportFields() {
        return this.meta.fields;
    }

    String getSupportLanguages() {
        return this.meta.languages.keySet().toString();
    }
}
