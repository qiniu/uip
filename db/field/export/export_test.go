package export

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qiniu/uip/db/field"
)

func TestFieldSelector(t *testing.T) {
	ast := assert.New(t)

	selector := ParseRule(DefaultRule)
	ast.Equal([]string{"country", "province", "city", "isp", "asn"}, selector.fields)
	fMap := []field.Pair{
		{field.Country, "c"},
		{field.Province, "p"},
		{field.City, "ct"},
		{field.ISP, "i"},
		{field.Asn, "a"},
		{field.Latitude, "l"},
	}
	data := map[string]string{
		"c":  "中国",
		"p":  "浙江",
		"ct": "杭州",
		"i":  "电信",
		"a":  "123",
		"cc": "ccc",
	}
	fMapFilter := selector.Select(fMap)
	ast.Equal([]field.Pair{
		{field.Country, "c"},
		{field.Province, "p"},
		{field.City, "ct"},
		{field.ISP, "i"},
		{field.Asn, "a"},
	}, fMapFilter)

	ast.Equal([]string{"中国", "浙江", "杭州", "电信", "123"}, selector.Export(fMapFilter, data))

	data = map[string]string{
		"c":  "日本",
		"p":  "东京都",
		"ct": "品川区",
		"i":  "WIDE Project",
		"a":  "123",
		"bb": "bbb",
	}
	ast.Equal([]string{"日本", "", "", "", ""}, selector.Export(fMapFilter, data))
}
