package ipdb

import (
	"github.com/qiniu/uip/db/inf"
	"io"
)

func pack(ipd *inf.IpData, writer io.Writer) error {
	meta := buildMeta(ipd.Version)
	meta.Fields = covertFields(ipd.Fields)
	w := NewWriter(*meta, nil)
	// writer insert ipraw from ipd
	for _, v := range ipd.Ips {
		w.insert(v.Cidr, v.FieldValues)
	}
	return w.Save(writer)
}
