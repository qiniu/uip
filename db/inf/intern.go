package inf

import (
	"github.com/qiniu/uip/db/field"
	"io"
	"net"
)

type IpData struct {
	Fields []string
	Ips    []IpRaw
}

type IpRaw struct {
	Cidr        *net.IPNet
	FieldValues []string
}

type NewDumper func([]byte) (Dump, error)

type Exporter interface {
	Fields() []string
	Select(fMap []field.Pair) []field.Pair
	Export(fieldMap []field.Pair, data map[string]string) []string
}

type Dump interface {
	// Dump does not support reenter
	Dump(exporter Exporter) (*IpData, error)
	VersionInfo() *VersionInfo
}

type NewPacker func() Pack

type Pack interface {
	Pack(ipd *IpData, v *VersionInfo, writer io.Writer) error
}

type Find func(ip net.IP) (cidr *net.IPNet, record map[string]string, err error)
