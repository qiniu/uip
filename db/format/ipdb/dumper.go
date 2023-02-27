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

	return export.BuildIPData(FieldsArray, exporter, db.version(),
		func(ip net.IP) (cidr *net.IPNet, record map[string]string, err error) {
			return db.findMap(ip, "CN")
		})
}
