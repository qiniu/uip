package main

import (
	"flag"
	"log"
	"strings"

	"github.com/qiniu/uip/db/convert"
	"github.com/qiniu/uip/db/field/operate"
	_ "github.com/qiniu/uip/db/format/awdb"
	_ "github.com/qiniu/uip/db/format/ipdb"
	_ "github.com/qiniu/uip/db/format/plain"
	"github.com/qiniu/uip/db/query"
)

func main() {
	input := flag.String("i", "", "input file")
	output := flag.String("o", "", "output file list,split by ','")
	rule := flag.String("r", "", "field select rule")
	line := flag.String("line", "", "isp line file")
	lineReplace := flag.Bool("line-replace", false, "replace line info")
	flag.Parse()
	d, err := convert.NewDumper(*input)
	if err != nil {
		log.Println(err)
		flag.PrintDefaults()
		return
	}
	ipData, ver, err := d.Dump(*rule)
	log.Println(len(ipData.Ips), ver, err)
	operate.ReplaceShortage(ipData)
	operate.TrimAsnIspDup(ipData)
	operate.MergeNearNetwork(ipData, ver)
	if *line != "" {
		lineQ, err := query.NewDb(*line)
		if err != nil {
			log.Println(err)
			flag.PrintDefaults()
			return
		}
		log.Println("line", lineQ.VersionInfo())
		if *lineReplace {
			operate.ReplaceLineByCidr(ipData, ver, query.GetIntern(lineQ))
		} else {
			operate.AttachLineByCidr(ipData, ver, query.GetIntern(lineQ))
		}
	}
	operate.AttachDistrict(ipData)

	outputs := strings.Split(*output, ",")
	for _, output := range outputs {
		err = convert.Pack(output, ipData, ver)
		if err != nil {
			log.Println(err)
			flag.PrintDefaults()
			return
		}
	}
}
