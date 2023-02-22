package query

import (
	"errors"
	"net"
	"path"
	"time"

	"github.com/qiniu/uip/db/inf"
	"github.com/qiniu/uip/db/query"
	"github.com/qiniu/uip/loader"
)

func NewQueryAgent(conf *Conf) (*Agent, error) {
	l, err := loader.NewLoader(conf.LoaderConf)
	if err != nil {
		return nil, err
	}
	var c = Agent{
		loader: l,
		conf:   conf,
	}
	err = c.sync()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// load db from loader by key
// if key is empty, return nil
func (p *Agent) load(key string, pdb **db, v6 bool) error {
	var ver string
	if *pdb != nil {
		ver = (*pdb).version
	}
	need, err := p.loader.NeedUpdate(key, ver)
	if err != nil {
		return err
	}
	if !need {
		return nil
	}
	if key == "" {
		return nil
	}
	kind := path.Ext(key)
	f := func(backup bool) (*query.Db, string, error) {
		content, version, err := p.loader.Load(key, backup)
		if err != nil {
			return nil, "", err
		}
		querier, err := query.NewDbFromBytes(kind, content)
		if err != nil {
			return nil, "", err
		}
		if v6 {
			err = querier.CheckV6()
		} else {
			err = querier.CheckV4()
		}
		if err != nil {
			return nil, "", err
		}
		return querier, version, nil
	}
	querier, version, err := f(false)
	if err != nil {
		querier, version, err = f(true)
		if err != nil {
			return err
		}
	}
	// ref go memory model, https://go.dev/ref/mem,
	// Otherwise, a read r of a memory location x that is not larger than a machine word must observe some write w
	// such that r does not happen before w and there is no write w' such that w happens before w' and w' happens before r.
	// That is, each read must observe a value written by a preceding or concurrent write.
	// write to *pdb is atomic, so no need to lock

	*pdb = &db{
		version: version,
		querier: querier,
	}
	return nil
}

func (p *Agent) sync() error {
	if p.conf.Ipv4Key != "" {
		err := p.load(p.conf.Ipv4Key, &p.ipv4db, false)
		if err != nil {
			return err
		}
	}

	if p.conf.Ipv6Key != "" {
		err := p.load(p.conf.Ipv6Key, &p.ipv6db, true)
		if err != nil {
			return err
		}
	}

	if p.conf.IpKey != "" {
		err := p.load(p.conf.IpKey, &p.ipdb, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Agent) Start() {
	p.loader.Start()
	go p.backgroundUpdate()
}

func (p *Agent) backgroundUpdate() {
	for {
		_ = p.sync()
		time.Sleep(12 * time.Hour)
	}
}

func (p *Agent) QueryStr(ip string) (*inf.IpInfo, error) {
	ipn := net.ParseIP(ip)
	if ipn == nil {
		return nil, errors.New("invalid ip " + ip)
	}
	return p.Query(ipn)
}

func (p *Agent) Query(ip net.IP) (*inf.IpInfo, error) {
	if ip.To4() != nil {
		if p.ipv4db != nil {
			return p.ipv4db.querier.Query(ip)
		} else if p.ipdb != nil {
			return p.ipdb.querier.Query(ip)
		}
	} else {
		if p.ipv6db != nil {
			return p.ipv6db.querier.Query(ip)
		} else if p.ipdb != nil {
			return p.ipdb.querier.Query(ip)
		}
	}
	return nil, errors.New("no db found")
}
