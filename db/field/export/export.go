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

func (e *Exporter) Select(fMap []field.Pair) []field.Pair {
	var ret = make([]field.Pair, 0, len(e.fields))
	for _, f := range e.fields {
		var p field.Pair
		for _, v := range fMap {
			if v.Ext == f {
				p = v
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

// todo support condition regression, country=中国*, country=in(中国,美国)
// ParseRule 构造字段选择器
// 用于选择输出的字段，支持简单的匹配逻辑
// 输入的参数格式为: <fields>[|<rule1>:<fields>|<rule2>:<fields>|<default fields>]
// @ <fields> - 字段列表，表示需要输出的字段，字段之间使用","分隔
// @ <rule> - 匹配规则，非必填项，以<key>=<value>表示，匹配上的话，则使用<rule>对应的<fields>，匹配优先级为<fields>的顺序
// @ sep - Select函数输出数据的分隔符，默认为","
// 举例
// country,province,city,isp,asn|country=!中国:country
// 针对国家区分IP库精度，若国家是"中国"，返回"国家,省份,城市,运营商,ASN"，若国家不是"中国"，返回"国家"

func ParseRule(rule string) *Exporter {
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

func BuildIPData(fieldArray []field.Pair, exporter inf.Exporter, version *inf.VersionInfo, find inf.Find) (*inf.IpData, error) {
	var ret inf.IpData
	var fieldsMap = fieldArray
	if exporter == nil {
		ret.Fields = field.ExtKeys(fieldArray)
	} else {
		ret.Fields = exporter.Fields()
		fieldsMap = exporter.Select(fieldArray)
	}
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

		if exporter != nil {
			values = exporter.Export(fieldsMap, info)
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
