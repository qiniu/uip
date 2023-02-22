package plain

import "github.com/qiniu/uip/db/format"

const ext = ".scan"

func init() {
	format.RegisterDumpFormat(ext, NewDumper)
	format.RegisterPackFormat(ext, NewPacker)
}
