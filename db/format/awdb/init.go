package awdb

import "github.com/qiniu/uip/db/format"

func init() {
	format.RegisterDumpFormat(".awdb", NewDumper)
	format.RegisterQueryFormat(".awdb", NewQuerier)
}
