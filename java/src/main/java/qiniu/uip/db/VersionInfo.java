package qiniu.uip.db;

import java.util.Arrays;

public final class VersionInfo {
    public static final int IPv4 = 1;
    public static final int IPv6 = 2;

    public int ipType;
    public int count;
    public long build;
    public String version;
    public String[] languages;
    public String[] extraInfo;

    public String toString() {
        return String.format("ipType: %d, count: %d, build: %d, version: %s, languages: %s, extraInfo: %s", ipType, count, build, version, Arrays.toString(languages), Arrays.toString(extraInfo));
    }

    public boolean hasIpV4() {
        return (ipType & IPv4) != 0;
    }

    public boolean hasIpV6() {
        return (ipType & IPv6) != 0;
    }

}
