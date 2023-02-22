package main

import (
	"log"
	"os"

	_ "github.com/qiniu/uip/db/format/awdb"
	_ "github.com/qiniu/uip/db/format/ipdb"
	"github.com/qiniu/uip/db/query"
)

func main() {
	f := os.Args[1]
	addr := os.Args[2]

	q, err := query.NewDb(f)
	if err != nil {
		log.Println(err)
		return
	}
	s := server{q}
	err = s.ListenAndServe(addr)
	if err != nil {
		log.Println(err)
		return
	}
}
