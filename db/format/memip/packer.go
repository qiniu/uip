package memip

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/qiniu/uip/db/inf"
	"github.com/qiniu/uip/ipnet"
)

func pack(ipd *inf.IpData, writer io.Writer) error {
	ipRanges := make([]RangeFields, 0, len(ipd.Ips))
	var prev RangeFields
	offsets := fieldOffsets(ipd.Fields)
	for _, ip := range ipd.Ips {
		values := covertValues(offsets, ip.FieldValues)
		if prev.r == nil {
			prev.r = ipnet.NewRange(ip.Cidr)
			prev.f = values
			continue
		}
		if reflect.DeepEqual(prev.f[1:], values[1:]) {
			if prev.r.JoinIPNet(ip.Cidr) {
				continue
			} else {
				ipRanges = append(ipRanges, prev)
				prev.r = ipnet.NewRange(ip.Cidr)
				prev.f = values
			}
		} else {
			ipRanges = append(ipRanges, prev)
			prev.r = ipnet.NewRange(ip.Cidr)
			prev.f = values
		}
	}
	ipRanges = append(ipRanges, prev)
	for _, rf := range ipRanges {
		view := strings.Join(rf.f, "-")
		if ipd.Version.IsIpV4() {
			fmt.Fprintf(writer, "%d %d %s\n", ipnet.IPv4ToUint32(rf.r.Start), ipnet.IPv4ToUint32(rf.r.End), view)
		} else {
			fmt.Fprintf(writer, "%s %s %s\n", rf.r.Start.String(), rf.r.End.String(), view)
		}
	}
	return nil
}

func fieldOffsets(fs []string) []int {
	ret := make([]int, len(fields))
	for i, f := range fields {
		ret[i] = -1
		for j, ff := range fs {
			if f == ff {
				ret[i] = j
				break
			}
		}
	}
	return ret
}

func covertValues(fieldsOff []int, values []string) []string {
	var newValues = make([]string, 6)
	for i, v := range fieldsOff {
		if v == -1 {
			newValues[i] = "默认"
		} else {
			nv := values[v]
			if nv == "" || strings.HasPrefix(nv, "保留") {
				nv = "默认"
			}
			// continent
			if i == 5 {
				switch nv {
				case "AS":
					nv = "亚洲"
				case "EU":
					nv = "欧洲"
				case "NA":
					nv = "北美洲"
				case "SA":
					nv = "南美洲"
				case "AF":
					nv = "非洲"
				case "OC":
					nv = "大洋洲"
				case "AN":
					nv = "南极洲"
				}
			}

			if i == 3 {
				if nv != "电信" && nv != "联通" && nv != "移动" && nv != "默认" {
					nv = "其他"
				}
			}
			newValues[i] = nv
		}
		pv := newValues[1]
		if pv == "香港" || pv == "澳门" || pv == "台湾" {
			newValues[2] = "默认"
		}

	}
	return newValues
}

type RangeFields struct {
	r *ipnet.Range
	f []string
}
