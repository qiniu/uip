package convert

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/qiniu/uip/db"
	"github.com/qiniu/uip/db/field/export"
	"github.com/qiniu/uip/db/format"
	"github.com/qiniu/uip/db/inf"
)

type Dumper struct {
	traverse inf.Dump
}

func NewDumper(ipFile string) (*Dumper, error) {
	data, err := os.ReadFile(ipFile)
	if err != nil {
		log.Println("Open ip file failed ", ipFile, err)
		return nil, err
	}
	extend := path.Ext(ipFile)
	create := format.GetDumpFormat(extend)
	if create == nil {
		log.Println("Unsupported ip file format ", ipFile, extend)
		return nil, db.ErrUnsupportedFormat
	}
	traverse, err := create(data)
	if err != nil {
		log.Println("New traverse failed ", ipFile, extend, err)
		return nil, err
	}

	return &Dumper{
		traverse: traverse,
	}, nil
}

func (d *Dumper) Dump(rule string) (*inf.IpData, *inf.VersionInfo, error) {
	v := d.traverse.VersionInfo()
	log.Println("version", v)
	if rule == "" {
		rule = export.DefaultRule
	}
	exp := export.ParseRule(rule)
	t0 := time.Now()
	all, err := d.traverse.Dump(exp)
	if err != nil {
		return nil, nil, err
	}
	log.Println("time elapse", time.Since(t0))
	return all, v, nil
}
