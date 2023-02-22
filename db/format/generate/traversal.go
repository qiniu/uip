package generate

// generate dummy ip infos for test

import (
	"github.com/qiniu/uip/db/inf"
	"time"
)

type traversal struct {
	Fields  []string
	Version inf.VersionInfo
}

func NewTraversal(lines uint32, fields []string, ipv inf.IpType) (tRet inf.Dump, err error) {
	t := &traversal{Fields: fields, Version: inf.VersionInfo{
		IpType:    ipv,
		Count:     lines,
		Build:     time.Now().Unix(),
		Version:   "generate-" + time.Now().Format("20060102"),
		Languages: nil,
		ExtraInfo: nil,
	}}

	return t, nil
}

func (t *traversal) Dump(exporter inf.Exporter) (*inf.IpData, error) {
	var ret inf.IpData
	return &ret, nil
}

func (t *traversal) VersionInfo() *inf.VersionInfo {
	return &t.Version
}
