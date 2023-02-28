package ipdb

import (
	db2 "github.com/qiniu/uip/db"
	"github.com/qiniu/uip/db/field"
	"github.com/qiniu/uip/db/inf"
	"net"
)

type querier struct {
	r          *reader
	queryCache map[int]*inf.IpInfo
}

func newQuerier(data []byte) (inf.Query, error) {
	db, err := newReaderFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &querier{db, make(map[int]*inf.IpInfo, 200)}, nil
}

func (q *querier) BuildCache(ipList []string) {
	for _, ip := range ipList {
		if ip == "" {
			continue
		}
		node, info, _, err := q.QueryInternal(net.ParseIP(ip), false)
		if err != nil || info == nil {
			continue
		}
		q.queryCache[node] = info
	}
}

func (q *querier) Query(ip net.IP) (*inf.IpInfo, int, error) {
	_, info, mask, err := q.QueryInternal(ip, true)
	return info, mask, err
}

func (q *querier) QueryInternal(ip net.IP, cache bool) (int, *inf.IpInfo, int, error) {
	node, body, mask, err := q.r.find0(ip)
	if cache {
		if v, ok := q.queryCache[node]; ok {
			return node, v, mask, nil
		}
	}

	off, ok := q.r.meta.Languages["CN"]
	if !ok {
		return 0, nil, 0, db2.ErrNoSupportLanguage
	}
	values, err := q.r.decodeInfo(body, off)
	if err != nil {
		return 0, nil, 0, err
	}

	ret := &inf.IpInfo{}
	for i, v := range q.r.meta.Fields {
		switch v {
		case FieldCountryName:
			ret.Country = values[i]
		case FieldRegionName:
			ret.Province = values[i]
		case FieldCityName:
			ret.City = values[i]
		case FieldISPDomain:
			ret.Isp = values[i]
		case field.Asn:
			ret.Asn = values[i]
		case FieldContinentCode:
			ret.Continent = values[i]
		case field.Line:
			ret.Line = values[i]
		case field.District:
			ret.District = values[i]
		}
	}
	return node, ret, mask, nil
}

func (q *querier) VersionInfo() *inf.VersionInfo {
	return q.r.version()
}
