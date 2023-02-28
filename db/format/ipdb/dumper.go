package ipdb

import (
	"github.com/qiniu/uip/db/field/export"
	"github.com/qiniu/uip/db/inf"
	"net"
)

func dump(data []byte, exporter inf.Exporter) (*inf.IpData, error) {
	db, err := newReaderFromBytes(data)
	if err != nil {
		return nil, err
	}
	ret, err := export.BuildIPData(CommonFieldsMap, db.meta.Fields, exporter, db.version(),
		func(ip net.IP) (cidr *net.IPNet, record map[string]string, err error) {
			ci, ret, err := db.findMap(ip, "CN")
			return ci, ret, err
		})
	if err != nil {
		return nil, err
	}
	ret.Version = db.version()
	return ret, err
}
