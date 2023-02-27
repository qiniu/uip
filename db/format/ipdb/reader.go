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

// Copy From https://github.com/ipipdotnet/ipdb-go
// modify by sjzar
// modify: Find / findMap return ipNet

import (
	"encoding/binary"
	"encoding/json"
	db2 "github.com/qiniu/uip/db"
	"github.com/qiniu/uip/db/inf"
	"net"
	"strings"
	"time"
	"unsafe"
)

type reader struct {
	fileSize  int
	nodeCount int
	v4offset  int

	meta            MetaData
	data            []byte
	continentOffset int
	ipV4Cache       []int
}

const cacheDepth = 12

func newReaderFromBytes(data []byte) (*reader, error) {
	var meta MetaData
	metaLength := int(binary.BigEndian.Uint32(data[0:4]))
	if len(data) < (4 + metaLength) {
		return nil, db2.ErrFileSize
	}
	if err := json.Unmarshal(data[4:4+metaLength], &meta); err != nil {
		return nil, err
	}
	if len(meta.Languages) == 0 || len(meta.Fields) == 0 {
		return nil, db2.ErrMetaData
	}
	if len(data) != (4 + metaLength + meta.TotalSize) {
		return nil, db2.ErrFileSize
	}

	db := &reader{
		fileSize:        len(data),
		nodeCount:       meta.NodeCount,
		meta:            meta,
		data:            data[4+metaLength:],
		continentOffset: -1,
	}
	//find offset of continent in meta.fields
	for i, v := range meta.Fields {
		if v == FieldContinentCode {
			db.continentOffset = i
			break
		}
	}

	if db.HasIPv4() {
		node := 0
		for i := 0; i < 96 && node < db.nodeCount; i++ {
			if i >= 80 {
				node = db.readNode(node, 1)
			} else {
				node = db.readNode(node, 0)
			}
		}
		db.v4offset = node
		db.initCache()
	}
	return db, nil
}

func indexToBytes(i int) []byte {
	return []byte{byte((i << (16 - cacheDepth)) >> 8), byte(0xFF & (i << (16 - cacheDepth)))}
}

func bytesToIndex(b []byte) int {
	return (0xFF&int(b[0]))<<8>>(16-cacheDepth) | (0xFF&int(b[1]))>>(16-cacheDepth)
}

func (db *reader) initCache() {
	db.ipV4Cache = make([]int, 1<<cacheDepth)
	//construct cache from binary trie tree for reduce read memory time
	for i := 0; i < len(db.ipV4Cache); i++ {
		b := indexToBytes(i)
		node, _ := db.readDepth(db.v4offset, cacheDepth, 0, b)
		db.ipV4Cache[i] = node
	}
	return
}

func (db *reader) findMap(addr net.IP, language string) (*net.IPNet, map[string]string, error) {
	ret, ipNet, err := db.find1(addr, language)
	if err != nil {
		return nil, nil, err
	}

	m := make(map[string]string)
	for i, v := range db.meta.Fields {
		m[v] = ret[i]
	}
	return ipNet, m, nil
}

func (db *reader) decodeInfo(body []byte, off int) ([]string, error) {
	str := (*string)(unsafe.Pointer(&body))
	tmp := strings.Split(*str, "\t")

	if (off + len(db.meta.Fields)) > len(tmp) {
		return nil, db2.ErrDatabaseError
	}
	ret := tmp[off : off+len(db.meta.Fields)]
	if db.continentOffset >= 0 {
		switch ret[db.continentOffset] {
		case "AS":
			ret[db.continentOffset] = "亚洲"
		case "EU":
			ret[db.continentOffset] = "欧洲"
		case "NA":
			ret[db.continentOffset] = "北美洲"
		case "SA":
			ret[db.continentOffset] = "南美洲"
		case "AF":
			ret[db.continentOffset] = "非洲"
		case "OC":
			ret[db.continentOffset] = "大洋洲"
		case "AN":
			ret[db.continentOffset] = "南极洲"
		}
	}
	return ret, nil
}

func (db *reader) find1(addr net.IP, language string) ([]string, *net.IPNet, error) {
	off, ok := db.meta.Languages[language]
	if !ok {
		return nil, nil, db2.ErrNoSupportLanguage
	}
	body, ipNet, err := db.find0(addr)
	if err != nil {
		return nil, nil, err
	}
	ret, err := db.decodeInfo(body, off)
	return ret, ipNet, err
}

func (db *reader) find0(ipv net.IP) ([]byte, *net.IPNet, error) {
	var bitCount int
	var ip net.IP
	if ip = ipv.To4(); ip != nil {
		if !db.HasIPv4() {
			return nil, nil, db2.ErrNoSupportIPv4
		}
		bitCount = 32
	} else if ip = ipv.To16(); ip != nil {
		if !db.HasIPv6() {
			return nil, nil, db2.ErrNoSupportIPv6
		}
		bitCount = 128
	} else {
		return nil, nil, db2.ErrIPFormat
	}

	node, mask, err := db.search(ip, bitCount)
	if err != nil || node < 0 {
		return nil, nil, err
	}

	cidrMask := net.CIDRMask(mask, len(ip)*8)
	ipNet := &net.IPNet{IP: ipv.Mask(cidrMask), Mask: cidrMask}

	body, err := db.resolve(node)
	if err != nil {
		return nil, nil, err
	}

	return body, ipNet, nil
}

func (db *reader) search(ip net.IP, bitCount int) (int, int, error) {
	var node int

	if bitCount == 32 {
		node = db.v4offset
	} else {
		node = 0
	}
	i := 0
	if db.ipV4Cache != nil && bitCount == 32 {
		node = db.ipV4Cache[bytesToIndex(ip)]
		if node > db.nodeCount {
			return node, i, nil
		}
		i = cacheDepth
	}

	node, i = db.readDepth(node, bitCount, i, ip)
	if node > db.nodeCount {
		return node, i, nil
	}

	return -1, 0, db2.ErrDataNotExists
}

func (db *reader) readDepth(node int, depth int, i int, ip []byte) (int, int) {
	for ; i < depth; i++ {
		if node >= db.nodeCount {
			break
		}
		node = db.readNode(node, ((0xFF&int(ip[i>>3]))>>uint(7-(i%8)))&1)
	}
	return node, i
}

func (db *reader) readNode(node, indexBit int) int {
	off := node*8 + indexBit*4
	return int(binary.BigEndian.Uint32(db.data[off : off+4]))
}

func (db *reader) resolve(node int) ([]byte, error) {
	resolved := node - db.nodeCount + db.nodeCount*8
	if resolved >= db.fileSize {
		return nil, db2.ErrDatabaseError
	}

	size := int(binary.BigEndian.Uint16(db.data[resolved : resolved+2]))
	if (resolved + 2 + size) > len(db.data) {
		return nil, db2.ErrDatabaseError
	}
	bytes := db.data[resolved+2 : resolved+2+size]

	return bytes, nil
}

func (db *reader) Build() time.Time {
	return time.Unix(db.meta.Build, 0).In(time.UTC)
}

func (db *reader) Languages() []string {
	ls := make([]string, 0, len(db.meta.Languages))
	for k := range db.meta.Languages {
		ls = append(ls, k)
	}
	return ls
}
func (db *reader) HasIPv4() bool {
	return (int(db.meta.IPVersion) & IPv4) == IPv4
}

func (db *reader) HasIPv6() bool {
	return (int(db.meta.IPVersion) & IPv6) == IPv6
}

func (db *reader) version() *inf.VersionInfo {
	var version string
	if db.meta.Version == "" {
		version = "ipip.net-" + db.Build().Format("2006-01-02")
	} else {
		version = db.meta.Version
	}

	return &inf.VersionInfo{
		IpType:    db.meta.IpType(),
		Count:     uint32(db.meta.NodeCount),
		Build:     db.meta.Build,
		Version:   version,
		Languages: db.Languages(),
		ExtraInfo: db.meta.Extra,
	}
}
