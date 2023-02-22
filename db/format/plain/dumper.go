package plain

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"strings"

	"github.com/qiniu/uip/db/field"
	"github.com/qiniu/uip/db/inf"
	"github.com/qiniu/uip/ipnet"
)

type traversal struct {
	lines   []string
	Fields  []string
	Version *inf.VersionInfo
}

func NewDumper(data []byte) (tRet inf.Dump, err error) {
	t := &traversal{}
	bio := bytes.NewReader(data)
	scanner := bufio.NewScanner(bio)
	lineOffset := 0
	for scanner.Scan() {
		if lineOffset == 0 {
			// ignore title
		} else if lineOffset == 1 {
			//	ignore prefix "#### "
			err := json.Unmarshal(scanner.Bytes()[5:], &t.Version)
			if err != nil {
				return nil, err
			}
		} else if lineOffset == 2 {
			//	ignore prefix "#### "
			b := scanner.Bytes()[5:]
			t.Fields = strings.Split(string(b), ",")
			t.lines = make([]string, 0, t.Version.Count+1)
		} else {
			t.lines = append(t.lines, scanner.Text())
		}
		lineOffset += 1
	}
	return t, nil
}

func fieldToPair(fields []string) []field.Pair {
	var ret []field.Pair
	for _, f := range fields {
		ret = append(ret, field.Pair{Ext: f, Intern: f})
	}
	return ret
}

func (t *traversal) parseLine(line string) (*net.IPNet, map[string]string, error) {
	var ret map[string]string
	ret = make(map[string]string, len(t.Fields))
	ar0 := strings.Split(line, "\t")
	if len(ar0) != 2 {
		return nil, nil, errors.New("invalid line: " + line)
	}
	_, cidr, err := net.ParseCIDR(ar0[0])
	if err != nil {
		return nil, nil, err
	}
	ar1 := strings.Split(ar0[1], ",")
	if len(ar1) != len(t.Fields) {
		return nil, nil, errors.New("invalid line: " + line)
	}
	for i, v := range ar1 {
		ret[t.Fields[i]] = v
	}
	return cidr, ret, nil
}

func (t *traversal) Dump(exporter inf.Exporter) (*inf.IpData, error) {
	var ret inf.IpData
	var fieldsMap []field.Pair
	if exporter == nil {
		ret.Fields = t.Fields
	} else {
		ret.Fields = exporter.Fields()
		fieldsMap = exporter.Select(fieldToPair(t.Fields))
	}
	ret.Ips = make([]inf.IpRaw, 0, t.Version.Count+1)
	var prevInfo []string
	for _, line := range t.lines {
		cidr, info, err := t.parseLine(line)
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
		prevInfo = values
	}

	return &ret, nil
}

func (t *traversal) VersionInfo() *inf.VersionInfo {
	return t.Version
}
