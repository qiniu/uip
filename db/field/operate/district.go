package operate

import (
	"github.com/qiniu/uip/db/field"
	"github.com/qiniu/uip/db/inf"
)

const ProvinceList = `上海市	上海
云南省	云南
北京市	北京
吉林省	吉林
四川省	四川
天津市	天津
安徽省	安徽
山东省	山东
山西省	山西
广东省	广东
江苏省	江苏
江西省	江西
河北省	河北
河南省	河南
浙江省	浙江
海南省	海南
湖北省	湖北
湖南省	湖南
甘肃省	甘肃
福建省	福建
贵州省	贵州
辽宁省	辽宁
重庆市	重庆
陕西省	陕西
青海省	青海
中国台湾	台湾
中国香港	香港
黑龙江省	黑龙江
西藏自治区	西藏
内蒙古自治区	内蒙古
宁夏回族自治区	宁夏
广西壮族自治区	广西
澳门特别行政区	澳门
新疆维吾尔自治区	新疆`

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
