package qiniu.uip.db.format.ipdb;

import qiniu.uip.db.IPFormatException;
import qiniu.uip.db.InvalidDatabaseException;
import qiniu.uip.db.NotFoundException;
import qiniu.uip.util.Json;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;


public final class Reader {

    MetaData meta;
    private int fileSize;
    private int nodeCount;
    private byte[] data;

    private int v4offset;

    public Reader(byte[] data) throws InvalidDatabaseException {
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
        }
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

    public String[] find(byte[] addr, String language) throws IPFormatException, InvalidDatabaseException {

        int off;
        try {
            off = this.meta.languages.get(language);
        } catch (NullPointerException e) {
            return null;
        }

        if (addr.length == 16) {
            if (!isIPv6()) {
                throw new IPFormatException("no support ipv6");
            }
        } else if (addr.length == 4) {
            if (!isIPv4()) {
                throw new IPFormatException("no support ipv4");
            }
        } else {
            throw new IPFormatException("ip format error");
        }

        // find node
        int node = 0;
        try {
            node = this.findNode(addr);
        } catch (NotFoundException nfe) {
            return null;
        }

        // resolve node
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

    private int findNode(byte[] binary) throws NotFoundException {

        int node = 0;

        final int bit = binary.length * 8;

        if (bit == 32) {
            node = this.v4offset;
        }

        for (int i = 0; i < bit; i++) {
            if (node > this.nodeCount) {
                break;
            }

            node = this.readNode(node, 1 & ((0xFF & binary[i / 8]) >> 7 - (i % 8)));
        }

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

    public boolean isIPv4() {
        return (this.meta.ip_version & 0x01) == 0x01;
    }

    public boolean isIPv6() {
        return (this.meta.ip_version & 0x02) == 0x02;
    }

    public int getBuildUTCTime() {
        return this.meta.build;
    }

    public String[] getSupportFields() {
        return this.meta.fields;
    }

    public String getSupportLanguages() {
        return this.meta.languages.keySet().toString();
    }
}