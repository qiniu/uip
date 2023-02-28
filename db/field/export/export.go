package export

import (
	"net"
	"reflect"
	"strings"

	"github.com/qiniu/uip/db/field"
	"github.com/qiniu/uip/db/inf"
	"github.com/qiniu/uip/ipnet"
)

const (
	SelectorGroupSep    = "|"
	SelectorFieldSep    = ","
	SelectorRuleSep     = ":"
	SelectorKeyValueSep = "="
	SelectorValueNot    = "!"
	SelectorValueOr     = "/"

	DefaultRule = "country,province,city,isp,asn,continent,district|country=!中国:country,continent|province=台湾/中国台湾:country,province,continent,district"
	LineRule    = "line"
)

type _Filter struct {
	keyOffset int
	condition struct {
		not    bool
		values []string
	}
	fieldsOffsets []int
}

type Exporter struct {
	fields  []string
	filters []*_Filter
}

func (e *Exporter) Fields() []string {
	return e.fields
}

func offsets(fields []string, keys []string) []int {
	var ret []int
	for _, k := range keys {
		for i, f := range fields {
			if f == k {
				ret = append(ret, i)
				break
			}
		}
	}
	return ret
}

func offset(fields []string, key string) int {
	for i, f := range fields {
		if f == key {
			return i
		}
	}
	return -1
}

func (e *Exporter) buildFilter(kc string, fields []string) *_Filter {
	var ret = &_Filter{
		fieldsOffsets: offsets(e.fields, fields),
	}
	ar := strings.Split(kc, SelectorKeyValueSep)
	if len(ar) != 2 {
		return nil
	}
	ret.keyOffset = offset(e.fields, ar[0])
	if ret.keyOffset == -1 {
		return nil
	}
	if strings.HasPrefix(ar[1], SelectorValueNot) {
		ret.condition.not = true
		ret.condition.values = strings.Split(ar[1][1:], SelectorValueOr)
	} else {
		ret.condition.values = strings.Split(ar[1], SelectorValueOr)
	}
	return ret
}

func (f *_Filter) match(value string) bool {
	in := field.InArray(value, f.condition.values)
	if f.condition.not {
		return !in
	}
	return in
}

func (f *_Filter) extract(data []string) []string {
	if f.keyOffset >= len(data) {
		return data
	} else {
		if f.match(data[f.keyOffset]) {
			var ret = make([]string, len(data))

			for _, offset := range f.fieldsOffsets {
				if offset < len(data) {
					ret[offset] = data[offset]
				}
			}
			return ret
		}
		return data
	}
}

func Select(fMap map[string]string, fields []string) []field.Pair {
	var ret = make([]field.Pair, 0, len(fields))
	for _, f := range fields {
		var p field.Pair
		for k, v := range fMap {
			if k == f {
				p.Intern = v
				p.Ext = f
				break
			}
		}
		ret = append(ret, p)
	}
	return ret
}

func (e *Exporter) remapKey(fieldMap []field.Pair, data map[string]string) []string {
	var ret = make([]string, len(e.fields))
	for k, v := range fieldMap {
		if f, ok := data[v.Intern]; ok {
			ret[k] = f
		}
	}
	return ret
}

func (e *Exporter) Export(fieldMap []field.Pair, data map[string]string) []string {
	ret := e.remapKey(fieldMap, data)
	if e.filters != nil {
		for _, filter := range e.filters {
			ret = filter.extract(ret)
		}
	}

	return ret
}

func ParseRule(rule string) inf.Exporter {
	if rule == "" {
		return nil
	}
	var ret = &Exporter{
		filters: make([]*_Filter, 0),
	}

	groups := strings.Split(rule, SelectorGroupSep)
	for i, group := range groups {
		fields := strings.Split(group, SelectorFieldSep)
		if i == 0 {
			ret.fields = fields
			continue
		}
		kc_v := strings.Split(fields[0], SelectorRuleSep)
		if len(kc_v) != 2 {
			continue
		}
		fields[0] = kc_v[1]
		var hasField = true
		for _, f := range fields {
			if !field.InArray(f, ret.fields) {
				hasField = false
				break
			}
		}
		if !hasField {
			continue
		}
		filter := ret.buildFilter(kc_v[0], fields)
		if filter == nil {
			continue
		}
		ret.filters = append(ret.filters, filter)
	}
	return ret
}

func BuildIPData(fieldMap map[string]string, fields []string, ept inf.Exporter, version *inf.VersionInfo, find inf.Find) (*inf.IpData, error) {
	var ret inf.IpData
	var fieldsMap = make([]field.Pair, 0, len(fields))
	if nil == ept {
		ret.Fields = field.ExtKeys(fieldMap, fields)
	} else {
		ret.Fields = ept.Fields()
	}
	fieldsMap = Select(fieldMap, ret.Fields)
	ret.Ips = make([]inf.IpRaw, 0, version.Count+1)
	var marker net.IP
	if version.IsIpV6() {
		marker = net.IPv6zero
	} else {
		marker = net.IPv4zero
	}

	var prevInfo []string
	for {
		cidr, info, err := find(marker)
		if err != nil {
			return nil, err
		}
		var values []string

		if ept != nil {
			values = ept.Export(fieldsMap, info)
		} else {
			for _, v := range ret.Fields {
				infK, ok := fieldMap[v]
				if !ok {
					values = append(values, "")
				} else if f, ok := info[infK]; ok {
					values = append(values, f)
				} else {
					values = append(values, "")
				}
			}
		}

		if reflect.DeepEqual(values, prevInfo) {
			values = prevInfo
		}

		ret.Ips = append(ret.Ips, inf.IpRaw{Cidr: cidr, FieldValues: values})
		last := ipnet.LastIP(cidr)
		if ipnet.IsAllMask(last) {
			break
		}
		marker = ipnet.NextStart(cidr)

		prevInfo = values
	}

	return &ret, nil
}
