package qiniu.uip.util;

import java.io.ByteArrayOutputStream;
import java.io.FileInputStream;
import java.io.IOException;

public final class File {
    private File() {
    }
    public static byte[] readAll(String path) throws IOException {
        FileInputStream in = new FileInputStream(path);
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        byte[] buffer = new byte[16384];
        int n;
        while ((n = in.read(buffer)) != -1) {
            out.write(buffer, 0, n);
        }
        byte[] b = out.toByteArray();
        return b;
    }


}
