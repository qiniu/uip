/*
 * Copyright (c) 2022 shenjunzheng@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ipdb

import (
	"github.com/qiniu/uip/db/field"
	"github.com/qiniu/uip/db/inf"
	"time"
)

// "country_name": "中国",
// "region_name": "浙江",
// "city_name": "",
// "isp_domain": "电信",
// "continent_code": "AP",
// "utc_offset": "UTC+8",
// "latitude": "29.19083",
// "longitude": "120.083656",
// "china_admin_code": "330000",
// "owner_domain": "",
// "timezone": "Asia/Shanghai",
// "idd_code": "86",
// "country_code": "CN",

const (
	FieldCountryName    = "country_name"
	FieldRegionName     = "region_name"
	FieldCityName       = "city_name"
	FieldISPDomain      = "isp_domain"
	FieldContinentCode  = "continent_code"
	FieldUTCOffset      = "utc_offset"
	FieldLatitude       = "latitude"
	FieldLongitude      = "longitude"
	FieldChinaAdminCode = "china_admin_code"
	FieldOwnerDomain    = "owner_domain"
	FieldTimezone       = "timezone"
	FieldIddCode        = "idd_code"
	FieldCountryCode    = "country_code"
	FieldIDC            = "idc"
	FieldBaseStation    = "base_station"
	FieldCountryCode3   = "country_code3"
	FieldEuropeanUnion  = "european_union"
	FieldCurrencyCode   = "currency_code"
	FieldCurrencyName   = "currency_name"
	FieldAnycast        = "anycast"
)

// FullFields 全字段列表
var FullFields = []string{
	FieldCountryName,
	FieldRegionName,
	FieldCityName,
	FieldISPDomain,
	FieldContinentCode,
	FieldUTCOffset,
	FieldLatitude,
	FieldLongitude,
	FieldChinaAdminCode,
	FieldOwnerDomain,
	FieldTimezone,
	FieldIddCode,
	FieldCountryCode,
	FieldIDC,
	FieldBaseStation,
	FieldCountryCode3,
	FieldEuropeanUnion,
	FieldCurrencyCode,
	FieldCurrencyName,
	FieldAnycast,
}

// CommonFieldsMap 公共字段映射
var CommonFieldsMap = map[string]string{
	field.Country:        FieldCountryName,
	field.Province:       FieldRegionName,
	field.City:           FieldCityName,
	field.ISP:            FieldISPDomain,
	field.Asn:            field.Asn,
	field.Continent:      FieldContinentCode,
	field.UTCOffset:      FieldUTCOffset,
	field.Latitude:       FieldLatitude,
	field.Longitude:      FieldLongitude,
	field.ChinaAdminCode: FieldChinaAdminCode,
	field.District:       field.District,
	field.Line:           field.Line,
}

const IPv4 = 0x01
const IPv6 = 0x02

type MetaData struct {
	Build     int64          `json:"build"`
	IPVersion uint16         `json:"ip_version"`
	Languages map[string]int `json:"languages"`
	NodeCount int            `json:"node_count"`
	TotalSize int            `json:"total_size"`
	Fields    []string       `json:"fields"`
	Version   string         `json:"version"`
	Extra     []string       `json:"extra"`
}

func covertFields(fields []string) []string {
	var newFields []string
	for _, f := range fields {
		if newF, ok := CommonFieldsMap[f]; ok {
			newFields = append(newFields, newF)
		} else {
			newFields = append(newFields, f)
		}
	}
	return newFields
}

func (m *MetaData) IpType() inf.IpType {
	if m.IPVersion == IPv4 {
		return inf.IpV4
	} else if m.IPVersion == IPv6 {
		return inf.IpV6
	} else {
		return inf.IpAll
	}
}

func ipVersion(v *inf.VersionInfo) uint16 {
	if v.IpType == inf.IpV4 {
		return IPv4
	} else if v.IpType == inf.IpV6 {
		return IPv6
	} else {
		return IPv4 | IPv6
	}
}

func buildMeta(v *inf.VersionInfo) *MetaData {
	meta := &MetaData{
		Build:     time.Now().Unix(),
		Version:   v.Version,
		IPVersion: ipVersion(v),
		Extra:     v.ExtraInfo,
	}
	return meta
}
