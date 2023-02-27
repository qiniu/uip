package plain

import "github.com/qiniu/uip/db/format"

const ext = ".scan"

func init() {
	format.RegisterDumpFormat(ext, dump)
	format.RegisterPackFormat(ext, pack)
}
