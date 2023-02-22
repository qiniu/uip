package operate

import (
	"reflect"

	"github.com/qiniu/uip/db/field"
	"github.com/qiniu/uip/db/field/data"
	"github.com/qiniu/uip/db/inf"
	"github.com/qiniu/uip/ipnet"
)

var mainIsp = []string{"电信", "联通", "移动", "铁通"}

func asnMatchIsp(s string, set map[string]map[string]bool) string {
	for k, v := range set {
		if _, ok := v[s]; ok {
			return k
		}
	}
	return ""
}

// TrimAsnIspDup some small isp use main isp bonenet whose asn is same as main isp
func TrimAsnIspDup(data *inf.IpData) {
	var ispOffset int
	var asnOffset int
	for i, f := range data.Fields {
		if f == field.ISP {
			ispOffset = i
		} else if f == field.Asn {
			asnOffset = i
		}
	}
	if ispOffset == 0 {
		return
	}
	var mainIspAsnMap = make(map[string]map[string]bool)
	for _, row := range data.Ips {
		isp := row.FieldValues[ispOffset]
		if asnOffset == 0 {
			continue
		}
		asn := row.FieldValues[asnOffset]
		if field.InArray(isp, mainIsp) {
			if _, ok := mainIspAsnMap[isp]; !ok {
				a := make(map[string]bool)
				a[asn] = true
				mainIspAsnMap[isp] = a
			} else {
				mainIspAsnMap[isp][asn] = true
			}
		}
	}

	if asnOffset == 0 {
		return
	}

	// replace some isp that asn appear in main isp
	for k, row := range data.Ips {
		isp := row.FieldValues[ispOffset]
		asn := row.FieldValues[asnOffset]
		if field.InArray(isp, mainIsp) {
			continue
		}

		if v := asnMatchIsp(asn, mainIspAsnMap); v != "" {
			data.Ips[k].FieldValues[ispOffset] = v
		}
	}
}

// ReplaceShortage should be the first step
func ReplaceShortage(ipdata *inf.IpData) {
	var offsets = make(map[int]map[string]string)
	for i, f := range ipdata.Fields {
		if v, ok := data.NormalizedMap[f]; ok {
			offsets[i] = v
		}
	}
	for k, row := range ipdata.Ips {
		for i, v := range row.FieldValues {
			if m, ok := offsets[i]; ok {
				if v, ok := m[v]; ok {
					ipdata.Ips[k].FieldValues[i] = v
				}
			}
		}
	}
}

type RangeFields struct {
	r *ipnet.Range
	f []string
}

func MergeNearNetwork(ipd *inf.IpData, v *inf.VersionInfo) {
	ipRanges := make([]RangeFields, 0, len(ipd.Ips))
	var prev RangeFields
	for _, ip := range ipd.Ips {
		if prev.r == nil {
			prev.r = ipnet.NewRange(ip.Cidr)
			prev.f = ip.FieldValues
			continue
		}
		if reflect.DeepEqual(prev.f, ip.FieldValues) {
			if prev.r.JoinIPNet(ip.Cidr) {
				continue
			} else {
				ipRanges = append(ipRanges, prev)
				prev.r = ipnet.NewRange(ip.Cidr)
				prev.f = ip.FieldValues
			}
		} else {
			ipRanges = append(ipRanges, prev)
			prev.r = ipnet.NewRange(ip.Cidr)
			prev.f = ip.FieldValues
		}
	}
	ipRanges = append(ipRanges, prev)
	ipd.Ips = make([]inf.IpRaw, 0, len(ipd.Ips))
	for _, ipRange := range ipRanges {
		subs := ipRange.r.IPNets()
		for _, sub := range subs {
			ipd.Ips = append(ipd.Ips, inf.IpRaw{
				Cidr:        sub,
				FieldValues: ipRange.f,
			})
		}
	}
	v.Count = uint32(len(ipd.Ips))
}
