package operate

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/qiniu/uip/db/field"
	"github.com/qiniu/uip/db/inf"
)

func AttachLineByCidr(data *inf.IpData, ver *inf.VersionInfo, line inf.Query) {
	lineVer := fmt.Sprintf("%s-%d", line.VersionInfo().Version, line.VersionInfo().Build)
	ver.ExtraInfo = append(ver.ExtraInfo, lineVer)
	var lineOffset int
	if lineOffset = field.Offset(field.Line, data.Fields); lineOffset == -1 {
		data.Fields = append(data.Fields, field.Line)
	}
	countryOffset := field.Offset(field.Country, data.Fields)
	if countryOffset == -1 {
		return
	}

	for k, row := range data.Ips {
		var ispline string
		if row.FieldValues[countryOffset] == "中国" {
			ispline = queryLine(row.Cidr, line)
		}
		if lineOffset == -1 {
			row.FieldValues = append(row.FieldValues, ispline)
		} else {
			row.FieldValues[lineOffset] = ispline
		}
		data.Ips[k] = row
	}
}

func AttachLineByAsn(data *inf.IpData, ver *inf.VersionInfo, asnLine map[string]string, asnLineVer string) {
	ver.ExtraInfo = append(ver.ExtraInfo, asnLineVer)
	var hasLine bool
	var lineOffset int
	var asnOffset int
	// if there is no asn field, return
	if !field.InArray(field.Asn, data.Fields) {
		return
	}
	if !field.InArray(field.Line, data.Fields) {
		data.Fields = append(data.Fields, field.Line)
		hasLine = true
		lineOffset = len(data.Fields) - 1
	} else {
		for i, f := range data.Fields {
			if f == field.Line {
				lineOffset = i
			} else if f == field.Asn {
				asnOffset = i
			}
		}
	}
	for k, row := range data.Ips {
		asn := row.FieldValues[asnOffset]
		if line, ok := asnLine[asn]; ok {
			if hasLine {
				row.FieldValues = append(row.FieldValues, line)
			} else if line == "" {
				row.FieldValues[lineOffset] = line
			}
			data.Ips[k] = row
		}
	}
}

func queryLine(ipnet *net.IPNet, line inf.Query) (ispline string) {
	ipStart := ipnet.IP
	l, err := line.Query(ipStart)
	if err != nil {
		log.Println(ipStart, err)
		return
	}
	// like 电信/联通/阿里云 just get the first part
	isp := strings.Split(l.Line, "/")
	ispline = isp[0]
	if ispline != "" && ispline != "电信" && ispline != "联通" && ispline != "移动" {
		if ispline == "铁通" {
			ispline = "移动"
		} else {
			// Replace all other line to 电信
			ispline = "电信"
		}
	}
	return
}

func ReplaceLineByCidr(data *inf.IpData, ver *inf.VersionInfo, line inf.Query) {
	lineVer := fmt.Sprintf("%s-%d", line.VersionInfo().Version, line.VersionInfo().Build)
	ver.ExtraInfo = append(ver.ExtraInfo, lineVer)
	ispOffset := field.Offset(field.ISP, data.Fields)
	if ispOffset == -1 {
		return
	}
	countryOffset := field.Offset(field.Country, data.Fields)
	if countryOffset == -1 {
		return
	}

	for k, row := range data.Ips {
		var ispline string
		if row.FieldValues[countryOffset] == "中国" {
			ispline = queryLine(row.Cidr, line)
		}
		if ispline != "" {
			row.FieldValues[ispOffset] = ispline
			data.Ips[k] = row
		}
	}
}
