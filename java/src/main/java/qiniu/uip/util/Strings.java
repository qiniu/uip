package qiniu.uip.util;

public final class Strings {
    private Strings() {
    }

    public static boolean isEmpty(String s) {
        return s == null || s.length() == 0;
    }

    public static String toHex(byte[] data) {
        StringBuilder sb = new StringBuilder();
        for (byte b : data) {
            String hex = Integer.toHexString(b & 0xFF);
            if (hex.length() == 1) {
                sb.append('0');
            }
            sb.append(hex);
        }
        return sb.toString();
    }
}
