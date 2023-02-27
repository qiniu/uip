package convert

import (
	"bytes"
	"io"
	"log"
	"os"
	"path"

	"github.com/qiniu/uip/db"
	"github.com/qiniu/uip/db/format"
	"github.com/qiniu/uip/db/inf"
)

func PackFile(ipFile string, ipData *inf.IpData) error {
	file, err := os.Create(ipFile)
	if err != nil {
		return err
	}
	defer file.Close()
	extend := path.Ext(ipFile)
	return PackWriter(extend, ipData, file)
}

func PackBytes(kind string, ipData *inf.IpData) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	err := PackWriter(kind, ipData, w)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func PackWriter(kind string, ipData *inf.IpData, writer io.Writer) error {
	create := format.GetPackFormat(kind)
	if create == nil {
		log.Println("Unsupported ip file format ", kind)
		return db.ErrUnsupportedFormat
	}

	return create(ipData, writer)
}
