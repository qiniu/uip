package qiniu.uip.db.format;

import qiniu.uip.db.InvalidDatabaseException;

@FunctionalInterface
public interface QueryBuilder {
    Query build(byte[] dat) throws InvalidDatabaseException;
}
