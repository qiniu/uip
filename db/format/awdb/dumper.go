package awdb

import (
	"github.com/qiniu/uip/db/field/export"
	"github.com/qiniu/uip/db/inf"
	"log"
	"net"
)

type traversal struct {
	Data []byte

	BuildEpoch   uint64
	DatabaseType string
	Description  map[string]string
	Languages    []string

	reader *Reader
}

func NewDumper(data []byte) (tRet inf.Dump, err error) {
	t := &traversal{Data: data}
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

func (t *traversal) check() error {
	if t.VersionInfo().HasIpV4() {
		r, m, err := t.Find(net.ParseIP("8.8.8.8"))
		if err != nil {
			return err
		}
		log.Println(r, m, err)
	}
	return nil
}

func (t *traversal) Dump(exporter inf.Exporter) (*inf.IpData, error) {
	return export.BuildIPData(FieldsArray, exporter, t.VersionInfo(), t.Find)
}

func (t *traversal) Find(ip net.IP) (*net.IPNet, map[string]string, error) {
	var record interface{}
	ipNet, _, err := t.reader.LookupNetwork(ip, &record)
	if err != nil {
		return nil, nil, err
	}

	data := make(map[string]string)
	for k, v := range record.(map[string]interface{}) {
		data[k] = string(v.([]byte))
	}
	return ipNet, data, nil
}

func (t *traversal) VersionInfo() *inf.VersionInfo {
	return t.reader.version()
}
