package query

import (
	"errors"
	"github.com/qiniu/uip/ipnet"
	"net"
	"os"
	"path"

	"github.com/qiniu/uip/db"
	"github.com/qiniu/uip/db/format"
	"github.com/qiniu/uip/db/inf"
)

type Db struct {
	q inf.Query
}

func NewDb(file string) (*Db, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	extend := path.Ext(file)
	return NewDbFromBytes(extend, b)
}

func NewDbFromBytes(kind string, b []byte) (*Db, error) {
	create := format.GetQueryFormat(kind)
	if create == nil {
		return nil, db.ErrUnsupportedFormat
	}
	q, err := create(b)
	if err != nil {
		return nil, err
	}
	return &Db{q}, nil
}

func (q *Db) QueryStr(ipStr string) (*inf.IpInfo, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, errors.New("invalid ip: " + ipStr)
	}
	return q.Query(ip)
}

func (q *Db) QueryU32(i uint32) (*inf.IpInfo, error) {
	ip := ipnet.Uint32ToIPv4(i)
	return q.Query(ip)
}

func (q *Db) CheckV4() error {
	google := net.ParseIP("8.8.8.8")
	info, err := q.Query(google)
	if err != nil {
		return err
	}
	if info.Country != "美国" {
		return db.ErrCheckFailed
	}
	dnspod := net.ParseIP("119.29.29.29")
	info, err = q.Query(dnspod)
	if err != nil {
		return err
	}
	if info.Country != "中国" {
		return db.ErrCheckFailed
	}
	return nil
}

func (q *Db) CheckV6() error {
	google := net.ParseIP("2001:4860:4860::8888")
	info, err := q.Query(google)
	if err != nil {
		return err
	}
	if info.Country != "美国" {
		return db.ErrCheckFailed
	}
	cloudFlare := net.ParseIP("2606:4700:4700::1111")
	info, err = q.Query(cloudFlare)
	if err != nil {
		return err
	}
	if info.Country != "美国" {
		return db.ErrCheckFailed
	}
	return nil
}

func (q *Db) Query(ip net.IP) (*inf.IpInfo, error) {
	info, err := q.q.Query(ip)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// version
func (q *Db) VersionInfo() *inf.VersionInfo {
	return q.q.VersionInfo()
}

// GetIntern only for internal use
func GetIntern(q *Db) inf.Query {
	return q.q
}
