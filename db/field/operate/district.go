package operate

import (
	"github.com/qiniu/uip/db/field"
	"github.com/qiniu/uip/db/inf"
)

// map province to district like "上海" => "华东"
var provinceDistrict = map[string]string{
	"黑龙江": "东北",
	"辽宁":  "东北",
	"吉林":  "东北",

	"新疆": "西北",
	"青海": "西北",
	"宁夏": "西北",
	"甘肃": "西北",
	"陕西": "西北",

	"内蒙古": "华北",
	"河北":  "华北",
	"山西":  "华北",
	"天津":  "华北",
	"北京":  "华北",

	"西藏": "西南",
	"四川": "西南",
	"重庆": "西南",
	"贵州": "西南",
	"云南": "西南",

	"湖南": "华中",
	"湖北": "华中",
	"河南": "华中",

	"山东": "华东",
	"江苏": "华东",
	"安徽": "华东",
	"浙江": "华东",
	"江西": "华东",
	"福建": "华东",
	"上海": "华东",

	"台湾": "华东",

	"广东": "华南",
	"广西": "华南",
	"海南": "华南",

	"香港": "华南",
	"澳门": "华南",
}

func AttachDistrict(data *inf.IpData) {
	districtOffset := field.Offset(field.District, data.Fields)
	if districtOffset == -1 {
		data.Fields = append(data.Fields, field.District)
	}

	provinceOffset := field.Offset(field.Province, data.Fields)
	if provinceOffset == -1 {
		return
	}

	for k, row := range data.Ips {
		var district string
		if p := row.FieldValues[provinceOffset]; p != "" {
			district = provinceDistrict[p]
		}
		if districtOffset == -1 {
			row.FieldValues = append(row.FieldValues, district)
		} else {
			row.FieldValues[districtOffset] = district
		}
		data.Ips[k] = row
	}
}
