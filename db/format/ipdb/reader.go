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
// modify: Find / FindMap return ipNet

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

	meta MetaData
	data []byte
}

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
		fileSize:  len(data),
		nodeCount: meta.NodeCount,
		meta:      meta,
		data:      data[4+metaLength:],
	}

	if db.v4offset == 0 {
		node := 0
		for i := 0; i < 96 && node < db.nodeCount; i++ {
			if i >= 80 {
				node = db.readNode(node, 1)
			} else {
				node = db.readNode(node, 0)
			}
		}
		db.v4offset = node
	}
	return db, nil
}

func (db *reader) FindMap(addr net.IP, language string) (*net.IPNet, map[string]string, error) {

	data, ipNet, err := db.find1(addr, language)
	if err != nil {
		return nil, nil, err
	}

	info := make(map[string]string, len(db.meta.Fields))
	for k, v := range data {
		info[db.meta.Fields[k]] = v
	}

	return ipNet, info, nil
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

	str := (*string)(unsafe.Pointer(&body))
	tmp := strings.Split(*str, "\t")

	if (off + len(db.meta.Fields)) > len(tmp) {
		return nil, nil, db2.ErrDatabaseError
	}

	return tmp[off : off+len(db.meta.Fields)], ipNet, nil
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

	var i = 0
	for ; i < bitCount; i++ {
		if node >= db.nodeCount {
			break
		}
		node = db.readNode(node, ((0xFF&int(ip[i>>3]))>>uint(7-(i%8)))&1)
	}

	if node > db.nodeCount {
		return node, i, nil
	}

	return -1, 0, db2.ErrDataNotExists
}

func (db *reader) readNode(node, index int) int {
	off := node*8 + index*4
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
