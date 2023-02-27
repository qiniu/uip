package inf

import (
	"github.com/qiniu/uip/db/field"
	"io"
	"net"
)

type IpData struct {
	Fields  []string
	Version *VersionInfo
	Ips     []IpRaw
}

type IpRaw struct {
	Cidr        *net.IPNet
	FieldValues []string
}

type Dump func([]byte, Exporter) (*IpData, error)

type Exporter interface {
	Fields() []string
	Select(fMap []field.Pair) []field.Pair
	Export(fieldMap []field.Pair, data map[string]string) []string
}

type Pack func(*IpData, io.Writer) error

type Find func(ip net.IP) (cidr *net.IPNet, record map[string]string, err error)
