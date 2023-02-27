package ipdb

import (
	"github.com/qiniu/uip/db/field"
	"github.com/qiniu/uip/db/inf"
	"net"
)

type Querier struct {
	r *reader
}

func newQuerier(data []byte) (inf.Query, error) {
	db, err := newReaderFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &Querier{db}, nil
}

func (q *Querier) Query(ip net.IP) (*inf.IpInfo, error) {
	info, cidr, err := q.r.find1(ip, "CN")
	if err != nil {
		return nil, err
	}

	ret := &inf.IpInfo{Cidr: *cidr}
	for k, v := range q.r.meta.Fields {
		switch v {
		case FieldCountryName:
			ret.Country = info[k]
		case FieldRegionName:
			ret.Province = info[k]
		case FieldCityName:
			ret.City = info[k]
		case FieldISPDomain:
			ret.Isp = info[k]
		case field.Asn:
			ret.Asn = info[k]
		case FieldContinentCode:
			ret.Continent = info[k]
		case field.Line:
			ret.Line = info[k]
		}
	}
	return ret, nil
}

func (q *Querier) VersionInfo() *inf.VersionInfo {
	return q.r.version()
}
