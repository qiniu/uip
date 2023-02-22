package awdb

// some code  from 		"github.com/dilfish/awdb-golang/awdb-golang"
// and modified to fit uip

//+--------------------------------+--------------------------------+
//|                          Node Chunk (20%)                       |
//+--------------------------------+--------------------------------+
//|                 Section Separator Size (16byte)                 |
//+--------------------------------+--------------------------------+
//|                          Data Chunk (79%)                       |
//+--------------------------------+--------------------------------+
//|    Metadata Start Marker []byte("\xAB\xCD\xEFipplus360.com")    |
//+--------------------------------+--------------------------------+
//|                             MetaData                            |
//+--------------------------------+--------------------------------+
//
//
//* 数据同样由3部分组成，分别是NodeChunk、DataChunk、MetaData
//* 分析解析库得知，NodeChunk设计有3种类型，每个节点由2个Record组成，RecordSize有24/28/32bit，
//* DataChunk中定义了多种数据类型，通过depth decode解析数据，DataChunk占据了数据库近80%的存储空间（~480M）
//* Metadata被放在了数据库的末尾，非json格式
//
//
//# Metadata 示例
// {
//    BinaryFormatMajorVersion:2
//    BinaryFormatMinorVersion:0
//    BuildEpoch:1631571147
//    DatabaseType:IP_city_single_WGS84_zh_net_awdb.awdb
//    Description:map[
//        cn:埃文科技
//        data_version:20210914
//        version:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//    ]
//    IPVersion:4
//    Languages:[cn]
//    NodeCount:17253031
//    RecordSize:32
//}
