package plain

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/qiniu/uip/db/inf"
)

func pack(ipd *inf.IpData, writer io.Writer) error {
	verStr, err := json.Marshal(ipd.Version)
	if err != nil {
		return err
	}
	fmt.Fprintf(writer, "#### scan dump %s ####\n", time.Now().String())
	fmt.Fprintf(writer, "#### %s\n", string(verStr))
	fmt.Fprintf(writer, "#### %s\n", strings.Join(ipd.Fields, ","))
	for _, v := range ipd.Ips {
		replaceComma(&v)
		_, err = fmt.Fprintf(writer, "%s\t%s\n", v.Cidr, strings.Join(v.FieldValues, ","))
		if err != nil {
			return err
		}
	}
	return nil
}

func replaceComma(s *inf.IpRaw) {
	for i, v := range s.FieldValues {
		s.FieldValues[i] = strings.Replace(v, ",", "_", -1)
	}
}
