package memip

import (
	"github.com/qiniu/uip/db/field"
	"github.com/qiniu/uip/db/format"
)

var fields = []string{field.City, field.Province, field.District, field.ISP, field.Country, field.Continent}

const ext = ".memip"

func init() {
	format.RegisterDumpFormat(ext, dump)
	format.RegisterPackFormat(ext, pack)
}
