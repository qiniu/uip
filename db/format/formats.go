package format

import (
	"github.com/qiniu/uip/db/inf"
)

func RegisterDumpFormat(ext string, create inf.Dump) {
	dumpFormats[ext] = create
}

func RegisterQueryFormat(ext string, create inf.NewQuery) {
	queryFormats[ext] = create
}

func RegisterPackFormat(ext string, create inf.Pack) {
	packFormats[ext] = create
}

func GetDumpFormat(ext string) inf.Dump {
	return dumpFormats[ext]
}

func GetQueryFormat(ext string) inf.NewQuery {
	return queryFormats[ext]
}

func GetPackFormat(ext string) inf.Pack {
	return packFormats[ext]
}

var dumpFormats = map[string]inf.Dump{}

var queryFormats = map[string]inf.NewQuery{}

var packFormats = map[string]inf.Pack{}
