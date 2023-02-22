package ipdb

import "github.com/qiniu/uip/db/format"

const ext = ".ipdb"

func init() {
	format.RegisterDumpFormat(ext, NewDumper)
	format.RegisterQueryFormat(ext, NewQuerier)
	format.RegisterPackFormat(ext, NewPacker)
}
