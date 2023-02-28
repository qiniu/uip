package main

import (
	"flag"
	"github.com/qiniu/uip/db/field/export"
	"github.com/qiniu/uip/db/inf"
	"log"
	"strings"
	"time"

	"github.com/qiniu/uip/db/convert"
	"github.com/qiniu/uip/db/field/operate"
	_ "github.com/qiniu/uip/db/format/awdb"
	_ "github.com/qiniu/uip/db/format/ipdb"
	_ "github.com/qiniu/uip/db/format/memip"
	_ "github.com/qiniu/uip/db/format/plain"
	"github.com/qiniu/uip/db/query"
)

func main() {
	input := flag.String("i", "", "input file")
	output := flag.String("o", "", "output file list,split by ','")
	rule := flag.String("r", "default", "field select rule")
	line := flag.String("line", "", "isp line file")
	lineReplace := flag.Bool("line-replace", false, "replace line info")
	flag.Parse()
	ops := operate.DefaultOperates
	if *line != "" {
		lineQ, err := query.NewDb(*line)
		if err != nil {
			log.Println(err)
			flag.PrintDefaults()
			return
		}
		log.Println("line", lineQ.VersionInfo())
		f := operate.AttachLineByCidr
		if *lineReplace {
			f = operate.ReplaceLineByCidr
		}
		ops = append(ops, func(data *inf.IpData) {
			f(data, query.GetIntern(lineQ))
		})
	}
	ops = append(ops, operate.AttachDistrict)
	r := *rule
	if r == "default" {
		r = export.DefaultRule
	}
	t0 := time.Now()
	ipData, err := convert.DumpFile(*input, r, ops)
	if err != nil {
		log.Println(err)
		flag.PrintDefaults()
		return
	}
	log.Println(len(ipData.Ips), ipData.Version, err, time.Since(t0))
	outputs := strings.Split(*output, ",")
	for _, output := range outputs {
		err = convert.PackFile(output, ipData)
		if err != nil {
			log.Println(err)
			flag.PrintDefaults()
			return
		}
	}
}
