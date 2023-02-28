package main

import (
	"encoding/json"
	"fmt"
	"github.com/qiniu/uip/db/query"
	"net/http"
)

type server struct {
	*query.Db
}

func (s *server) ListenAndServe(addr string) error {
	m := http.NewServeMux()
	m.HandleFunc("/version", func(writer http.ResponseWriter, request *http.Request) {
		b, _ := json.Marshal(s.VersionInfo())
		fmt.Println(string(b))
		writer.Write(b)
	})
	m.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		ip := request.URL.Path[1:]
		i, _, err := s.QueryStr(ip)
		if err != nil {
			writer.WriteHeader(404)
			return
		}
		b, _ := json.Marshal(i)
		writer.Write(b)
	})

	serv := http.Server{
		Addr:    addr,
		Handler: m,
	}
	return serv.ListenAndServe()
}
