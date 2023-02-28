package qiniu.uip.util;

import sun.net.util.IPAddressUtil;

public final class IPAddress {
    private IPAddress() {
    }

    public static byte[] textToNumericFormatV4(String src) {
        return IPAddressUtil.textToNumericFormatV4(src);
    }

    public static byte[] textToNumericFormatV6(String src) {
        return IPAddressUtil.textToNumericFormatV6(src);
    }
}
