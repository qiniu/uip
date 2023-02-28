package qiniu.uip.db.format;

import java.util.HashMap;
import java.util.Map;

public final class Formats {
    private static Map<String, QueryBuilder> queryFormats = new HashMap<>();

    private Formats() {
    }

    public static void registerQueryFormat(String ext, QueryBuilder q) {
        queryFormats.put(ext, q);
    }

    public static QueryBuilder getQueryFormat(String ext) {
        return queryFormats.get(ext);
    }
}
