package awdb

import (
	"log"
	"net"

	"github.com/qiniu/uip/db/inf"
)

type querier struct {
	Data []byte

	BuildEpoch   uint64
	DatabaseType string
	Description  map[string]string
	Languages    []string

	reader *Reader
}

func newQuerier(data []byte) (tRet inf.Query, err error) {
	t := &querier{Data: data}
	t.reader, err = FromBytes(data)
	if err != nil {
		return nil, err
	}

	err = t.check()
	if err != nil {
		return nil, err
	}
	return t, err
}

func (q *querier) Query(ip net.IP) (*inf.IpInfo, int, error) {
	cidr, info, err := q.Find(ip)
	if err != nil {
		return nil, 0, err
	}
	ones, _ := cidr.Mask.Size()
	return &inf.IpInfo{
		Country:   info[FieldCountry],
		Province:  info[FieldProvince],
		City:      info[FieldCity],
		Isp:       info[FieldISP],
		Asn:       info[FieldASNumber],
		Continent: info[FieldContinent],
		Line:      info[FieldLine],
	}, ones, nil
}

func (t *querier) check() error {
	if t.VersionInfo().HasIpV4() {
		r, m, err := t.Find(net.ParseIP("8.8.8.8"))
		if err != nil {
			return err
		}
		log.Println(r, m, err)
	}
	return nil
}

func (q *querier) BuildCache(ipList []string) {
}

func (q *querier) Find(ip net.IP) (*net.IPNet, map[string]string, error) {
	var record interface{}
	ipNet, _, err := q.reader.LookupNetwork(ip, &record)
	if err != nil {
		return nil, nil, err
	}

	data := make(map[string]string)
	for k, v := range record.(map[string]interface{}) {
		data[k] = string(v.([]byte))
	}
	return ipNet, data, nil
}

func (q *querier) VersionInfo() *inf.VersionInfo {
	return q.reader.version()
}
