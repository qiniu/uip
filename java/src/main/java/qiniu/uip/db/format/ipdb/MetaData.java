package qiniu.uip.db.format.ipdb;

import java.util.Map;

@SuppressWarnings("checkstyle:MemberName")
public final class MetaData {
    public int build;
    public int ip_version;
    public int node_count;
    public Map<String, Integer> languages;
    public String[] fields;
    public int total_size;

    public String version;
    public String[] extra;
}
