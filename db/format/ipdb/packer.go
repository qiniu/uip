package ipdb

import (
	"github.com/qiniu/uip/db/inf"
	"io"
)

type Packer struct {
	w *Writer
}

func NewPacker() inf.Pack {
	return &Packer{}
}

func (p *Packer) Pack(ipd *inf.IpData, v *inf.VersionInfo, writer io.Writer) error {
	meta := buildMeta(v)
	meta.Fields = covertFields(ipd.Fields)
	p.w = NewWriter(*meta, nil)
	// writer insert ipraw from ipd
	for _, v := range ipd.Ips {
		p.w.insert(v.Cidr, v.FieldValues)
	}
	return p.w.Save(writer)
}
