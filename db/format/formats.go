package format

import (
	"github.com/qiniu/uip/db/inf"
)

func RegisterDumpFormat(ext string, create inf.NewDumper) {
	dumpFormats[ext] = create
}

func RegisterQueryFormat(ext string, create inf.NewQuerier) {
	queryFormats[ext] = create
}

func RegisterPackFormat(ext string, create inf.NewPacker) {
	packFormats[ext] = create
}

func GetDumpFormat(ext string) inf.NewDumper {
	return dumpFormats[ext]
}

func GetQueryFormat(ext string) inf.NewQuerier {
	return queryFormats[ext]
}

func GetQueryFormatByDetect(b []byte) inf.NewQuerier {
	for k, detect := range queryRawFormats {
		if detect(b) {
			return queryFormats[k]
		}
	}
	return nil
}

func GetPackFormat(ext string) inf.NewPacker {
	return packFormats[ext]
}

type DetectFormat func([]byte) bool

var dumpFormats = map[string]inf.NewDumper{}

var queryFormats = map[string]inf.NewQuerier{}
var queryRawFormats = map[string]DetectFormat{}

var packFormats = map[string]inf.NewPacker{}
