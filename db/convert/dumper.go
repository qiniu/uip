package convert

import (
	"github.com/qiniu/uip/db"
	"github.com/qiniu/uip/db/field/export"
	"github.com/qiniu/uip/db/field/operate"
	"github.com/qiniu/uip/db/format"
	"github.com/qiniu/uip/db/inf"
	"log"
	"os"
	"path"
)

func DumpFile(ipFile, rule string, ops []operate.Operate) (*inf.IpData, error) {
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
	if rule == "" {
		rule = export.DefaultRule
	}
	ipData, err := create(data, export.ParseRule(rule))
	if err != nil {
		log.Println("Dump ip file failed ", ipFile, err)
		return nil, err
	}

	for _, op := range ops {
		op(ipData)
	}
	return ipData, nil
}
