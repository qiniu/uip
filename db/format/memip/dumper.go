package memip

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/qiniu/uip/db/inf"
	"github.com/qiniu/uip/ipnet"
)

func dump(data []byte, exporter inf.Exporter) (*inf.IpData, error) {
	t, err := newDumper(data)
	if err != nil {
		return nil, err
	}
	i, err := t.Dump(exporter)
	if err != nil {
		return nil, err
	}
	i.Version = t.Version
	return i, nil
}

type traversal struct {
	lines   []string
	Fields  []string
	Version *inf.VersionInfo
}

func newDumper(data []byte) (tRet *traversal, err error) {
	tm := time.Now()
	t := &traversal{
		//{城市}-{省份}-{大区}-{运营商}-{国家}-{大洲}
		Fields: fields,
		Version: &inf.VersionInfo{
			IpType:    0,
			Count:     0,
			Build:     tm.Unix(),
			Version:   tm.Format("20060102"),
			Languages: nil,
			ExtraInfo: nil,
		},
	}
	bio := bytes.NewReader(data)
	scanner := bufio.NewScanner(bio)
	for scanner.Scan() {
		t.lines = append(t.lines, scanner.Text())
	}
	return t, nil
}

func splitView(view string) []string {
	return strings.Split(view, "-")
}

func parseV6Range(ar []string) *ipnet.Range {
	start := net.ParseIP(ar[0])
	if start == nil {
		return nil
	}
	end := net.ParseIP(ar[1])
	if end == nil {
		return nil
	}
	return &ipnet.Range{Start: start, End: end}
}

func parseV4Range(ar []string) *ipnet.Range {
	start, err := strconv.ParseUint(ar[0], 10, 32)
	if err != nil {
		return nil
	}
	end, err := strconv.ParseUint(ar[1], 10, 32)
	if err != nil {
		return nil
	}
	return &ipnet.Range{Start: ipnet.Uint32ToIPv4(uint32(start)), End: ipnet.Uint32ToIPv4(uint32(end))}
}

func parseLine(line string) (*ipnet.Range, []string, error) {
	array := strings.Split(line, " ")
	if len(array) != 3 {
		return nil, nil, errors.New("bad format:" + line)
	}
	var rg *ipnet.Range
	if strings.Contains(array[0], ":") {
		rg = parseV6Range(array[0:2])
	} else {
		rg = parseV4Range(array[0:2])
	}

	view := array[2]
	return rg, splitView(view), nil
}

type row struct {
	rg   *ipnet.Range
	view []string
}

func (t *traversal) ipVer() inf.IpType {
	if strings.Contains(t.lines[0], ":") {
		return inf.IpV6
	}
	return inf.IpV4
}

func (t *traversal) Dump(_ inf.Exporter) (*inf.IpData, error) {
	var rows []row

	var prevInfo []string
	for _, line := range t.lines {
		cidr, values, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		if reflect.DeepEqual(values, prevInfo) {
			values = prevInfo
		}
		rows = append(rows, row{rg: cidr, view: values})
		if ipnet.IsAllMask(cidr.End) {
			break
		}
		prevInfo = values
	}
	var ipd = inf.IpData{
		Fields: t.Fields,
	}
	ipd.Ips = make([]inf.IpRaw, 0, len(ipd.Ips))
	for _, ipRange := range rows {
		subs := ipRange.rg.IPNets()
		for _, sub := range subs {
			ipd.Ips = append(ipd.Ips, inf.IpRaw{
				Cidr:        sub,
				FieldValues: ipRange.view,
			})
		}
	}
	t.Version.Count = uint32(len(ipd.Ips))
	t.Version.IpType = t.ipVer()
	return &ipd, nil
}

func (t *traversal) VersionInfo() *inf.VersionInfo {
	return t.Version
}
