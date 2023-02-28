package qiniu.uip.util;

import com.google.gson.Gson;

public final class Json {
    private Json() {
    }

    public static <T> T parseObject(String json, Class<T> clazz) {
        return (new Gson()).fromJson(json, clazz);
    }
}
