package inf

import (
	"encoding/json"
	"net"
)

type IpInfo struct {
	Country   string `json:"country,omitempty"`
	District  string `json:"district,omitempty"`
	Province  string `json:"province,omitempty"`
	City      string `json:"city,omitempty"`
	Asn       string `json:"asn,omitempty"`
	Isp       string `json:"isp,omitempty"`
	Continent string `json:"continent,omitempty"`

	//line字段，是指 IP是通过哪个运营商连接上全国骨干网的， 主要是通过看AS的连接关系，
	//比如AS55990与4811,4816,23724,58466,138950,9808,24400等AS连接，line就是 电信/联通/移动/华为；
	Line string `json:"line,omitempty"`
}

type NewQuery func([]byte) (Query, error)

type Query interface {
	Query(ip net.IP) (*IpInfo, int, error)
	BuildCache(ipList []string)
	VersionInfo() *VersionInfo
}

type VersionInfo struct {
	IpType    IpType
	Count     uint32
	Build     int64
	Version   string
	Languages []string
	ExtraInfo []string
}

type IpType int

const (
	IpV4  IpType = 1
	IpV6  IpType = 2
	IpAll IpType = 3
)

func (v *VersionInfo) HasIpV4() bool {
	return v.IpType&IpV4 != 0
}

func (v *VersionInfo) HasIpV6() bool {
	return v.IpType&IpV6 != 0
}

func (v *VersionInfo) IsIpAll() bool {
	return v.IpType == IpAll
}

func (v *VersionInfo) IsIpV4() bool {
	return v.IpType == IpV4
}

func (v *VersionInfo) IsIpV6() bool {
	return v.IpType == IpV6
}

func (v *VersionInfo) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}
