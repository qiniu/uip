universal IP query/dump/pack library
===================================
Support for IPv4, IPv6
Support Format: awdb, ipdb, plain text
Support Language: Go, Java

Convert Format
--------------
the source format need support dump, the target format need support pack
```
convert export rule: : <default fields>[|<rule1>:<only fields>|<rule2>:<only fields>]
if rule start with "!", it means not match
the rule condition can be use / to split, it means or
example: "country,province,city,isp,asn,continent,district|country=!中国:country,continent|province=台湾/中国台湾:country,province,continent,district"
```

```
convert example:
./cv -i ~/practice/dataset/ipv4.awdb -o new.scan[,new.ipdb]
the command can output multi fomat, the format is split by ",", the scan is text format for diff version update. 
```
