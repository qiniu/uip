package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/qiniu/uip/db/convert"
	_ "github.com/qiniu/uip/db/format/awdb"
	_ "github.com/qiniu/uip/db/format/ipdb"
	_ "github.com/qiniu/uip/db/format/plain"
)

// this package is used to check the difference between two ip database by comparing the scan result

func main() {
	left := os.Args[1]
	right := os.Args[2]
	dif := os.Args[3]
	dumpScan(left)
	dumpScan(right)
	cmd := exec.Command("diff", left+".scan", right+".scan")
	f, err := os.Create(dif)
	if err != nil {
		log.Println(err)
		flag.PrintDefaults()
		return
	}
	cmd.Stdout = f
	cmd.Run()
	f.Close()
}

func dumpScan(file string) {
	leftD, err := convert.NewDumper(file)
	if err != nil {
		log.Println(err)
		flag.PrintDefaults()
		return
	}
	data, ver, err := leftD.Dump("")

	err = convert.Pack(file+".scan", data, ver)
	if err != nil {
		log.Println(err)
		flag.PrintDefaults()
		return
	}

}
