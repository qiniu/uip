package awdb

import "github.com/qiniu/uip/db/format"

func init() {
	format.RegisterDumpFormat(".awdb", dump)
	format.RegisterQueryFormat(".awdb", newQuerier)
}
