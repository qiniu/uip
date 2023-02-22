package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/qiniu/uip/db/format/awdb"
	_ "github.com/qiniu/uip/db/format/ipdb"
	"github.com/qiniu/uip/db/query"
)

func main() {
	f := os.Args[1]
	ip := os.Args[2]

	q, err := query.NewDb(f)
	if err != nil {
		log.Println(err)
		return
	}
	i, err := q.QueryStr(ip)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("%+v\n", i)
}
