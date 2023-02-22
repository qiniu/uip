package convert

import (
	"bytes"
	"log"
	"os"
	"path"

	"github.com/qiniu/uip/db"
	"github.com/qiniu/uip/db/format"
	"github.com/qiniu/uip/db/inf"
)

func Pack(ipFile string, ipData *inf.IpData, info *inf.VersionInfo) error {
	file, err := os.Create(ipFile)
	if err != nil {
		return err
	}
	defer file.Close()
	extend := path.Ext(ipFile)
	create := format.GetPackFormat(extend)
	if create == nil {
		log.Println("Unsupported ip file format ", ipFile, extend)
		return db.ErrUnsupportedFormat
	}
	packer := create()
	log.Println("packing", ipFile, info)
	err = packer.Pack(ipData, info, file)
	log.Println("pack done", ipFile, err)
	return err
}

func PackBytes(kind string, ipData *inf.IpData, info *inf.VersionInfo) ([]byte, error) {
	create := format.GetPackFormat(kind)
	if create == nil {
		log.Println("Unsupported ip file format ", kind)
		return nil, db.ErrUnsupportedFormat
	}
	packer := create()
	w := bytes.NewBuffer(nil)
	err := packer.Pack(ipData, info, w)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
