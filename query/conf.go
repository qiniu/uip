package query

import (
	"encoding/json"
	"os"

	"github.com/qiniu/uip/db/query"
	"github.com/qiniu/uip/loader"
)

type db struct {
	version string
	querier *query.Db
}

type Agent struct {
	loader *loader.Loader
	ipdb   *db
	ipv4db *db
	ipv6db *db
	conf   *Conf
}

type Conf struct {
	LoaderConf *loader.Conf `json:"loader_conf,omitempty"`
	Ipv6Key    string       `json:"ipv6_key,omitempty"`
	Ipv4Key    string       `json:"ipv4_key,omitempty"`
	IpKey      string       `json:"ip_key,omitempty"`
}

func LoadConf(conf string) (*Conf, error) {
	c := &Conf{}
	f, err := os.Open(conf)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
