package ipdb

import (
	"github.com/qiniu/uip/db/field/export"
	"github.com/qiniu/uip/db/inf"
	"log"
	"net"
)

type traversal struct {
	r *reader
}

func NewDumper(data []byte) (tRet inf.Dump, err error) {
	db, err := newReaderFromBytes(data)
	if err != nil {
		return nil, err
	}
	t := &traversal{db}
	t.check()

	return t, nil
}

func (t *traversal) check() error {
	if t.VersionInfo().HasIpV4() {
		r, m, err := t.r.FindMap(net.ParseIP("8.8.8.8"), "CN")
		if err != nil {
			return err
		}
		log.Println(r, m, err)
	}
	return nil
}

func (t *traversal) Dump(exporter inf.Exporter) (*inf.IpData, error) {
	return export.BuildIPData(FieldsArray, exporter, t.VersionInfo(),
		func(ip net.IP) (cidr *net.IPNet, record map[string]string, err error) {
			return t.r.FindMap(ip, "CN")
		})
}

func (t *traversal) VersionInfo() *inf.VersionInfo {
	return t.r.version()
}
