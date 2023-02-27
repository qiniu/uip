package ipdb

import "github.com/qiniu/uip/db/format"

const ext = ".ipdb"

func init() {
	format.RegisterDumpFormat(ext, dump)
	format.RegisterQueryFormat(ext, newQuerier)
	format.RegisterPackFormat(ext, pack)
}
