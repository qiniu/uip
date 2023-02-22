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
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/qiniu/uip/db"
	"io"
	"log"
	"net"
	"strings"
)

const (
	FieldsSep = "\t"
)

// Writer IPDB 写入工具
type Writer struct {
	Meta      MetaData
	node      [][2]int
	dataHash  map[string]int
	dataChunk *bytes.Buffer
}

// NewWriter 初始化 IPDB 写入实例
func NewWriter(meta MetaData, languages map[string]int) *Writer {
	if len(languages) == 0 {
		meta.Languages = map[string]int{"CN": 0}
	} else {
		meta.Languages = languages
	}

	return &Writer{
		Meta:      meta,
		node:      [][2]int{{}},
		dataChunk: &bytes.Buffer{},
		dataHash:  make(map[string]int),
	}
}

// Save 保存数据
func (p *Writer) Save(w io.Writer) error {

	// Node Chunk
	nodeChunk := &bytes.Buffer{}
	p.Meta.NodeCount = len(p.node)
	for i := 0; i < p.Meta.NodeCount; i++ {
		for j := 0; j < 2; j++ {
			// 小于0: 数据记录，设置为NodeLength+Data偏移量
			// 等于0: 空值，设置为NodeLength
			// 大于0: Node跳转记录，不做调整
			if p.node[i][j] <= 0 {
				p.node[i][j] = p.Meta.NodeCount - p.node[i][j]
			}
			nodeChunk.Write(IntToBinaryBE(p.node[i][j], 32))
		}
	}
	// loopBack node
	nodeChunk.Write(IntToBinaryBE(p.Meta.NodeCount, 32))
	nodeChunk.Write(IntToBinaryBE(p.Meta.NodeCount, 32))

	// MetaData Chunk
	metaDataChunk := &bytes.Buffer{}
	p.Meta.TotalSize = nodeChunk.Len() + p.dataChunk.Len()
	metaData, err := json.Marshal(p.Meta)
	if err != nil {
		return err
	}
	metaDataChunk.Write(IntToBinaryBE(len(metaData), 32))
	metaDataChunk.Write(metaData)

	// Result
	if _, err := metaDataChunk.WriteTo(w); err != nil {
		return err
	}
	if _, err := nodeChunk.WriteTo(w); err != nil {
		return err
	}
	if _, err := p.dataChunk.WriteTo(w); err != nil {
		return err
	}

	return nil
}

// insert 插入数据
func (p *Writer) insert(ipNet *net.IPNet, values []string) error {
	mask, _ := ipNet.Mask.Size()
	node, index, ok := p.Nodes(ipNet.IP, mask)
	if !ok {
		log.Printf("load cidr failed cidr(%s) data(%s) node(%d) index(%d) preview data(%s)\n", ipNet, values, node, index, p.resolve(-p.node[node][index]))
		return db.ErrInvalidCIDR
	}
	if p.node[node][index] > 0 {
		log.Printf("cidr conflict %s %s\n", ipNet, values)
		return db.ErrCIDRConflict
	}
	offset := p.Fields(values)
	p.node[node][index] = -offset
	return nil
}

// resolve 解析数据
func (p *Writer) resolve(offset int) string {
	offset -= 8
	data := p.dataChunk.Bytes()
	size := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	if (offset + 2 + size) > len(data) {
		return ""
	}
	return string(data[offset+2 : offset+2+size])
}

// Nodes 获取CIDR地址所在节点和index
// 将补全Node中间链路，如果中间链路已经有数据，将无法写入新数据
func (p *Writer) Nodes(ip net.IP, mask int) (node, index int, ok bool) {
	// 如果传入的是IPv4，子网掩码增加96位( len(IPv6)-len(IPv4) )
	// 统一扩展为IPv6的子网掩码进行处理
	maxMask := mask - 1
	if ip.To4() != nil {
		if maxMask < 32 {
			maxMask += 96
		}
		if len(ip) == net.IPv4len {
			ip = ip.To16()
		}
	}
	for i := 0; i < maxMask; i++ {
		index = ((0xFF & int(ip[i>>3])) >> uint(7-(i%8))) & 1
		if p.node[node][index] == 0 {
			p.node = append(p.node, [2]int{})
			p.node[node][index] = len(p.node) - 1
		}
		if p.node[node][index] < 0 {
			return node, index, false
		}
		node = p.node[node][index]
	}
	return node, ((0xFF & int(ip[maxMask>>3])) >> uint(7-(maxMask%8))) & 1, true
}

// Fields 保存数据并返回数据的偏移量
// 相同的数据仅保存一份
// 数据格式 2 byte length + n byte data
func (p *Writer) Fields(fields []string) int {
	data := strings.Join(fields, FieldsSep)
	if _, ok := p.dataHash[data]; !ok {
		_data := []byte(data)
		// +8 是由于 loopBack node 占用了8byte
		p.dataHash[data] = p.dataChunk.Len() + 8
		p.dataChunk.Write(IntToBinaryBE(len(_data), 16))
		p.dataChunk.Write(_data)
	}
	return p.dataHash[data]
}

// IntToBinaryBE 将int转换为 binary big endian
func IntToBinaryBE(num, length int) []byte {
	switch length {
	case 16:
		_num := uint16(num)
		return []byte{byte(_num >> 8), byte(_num)}
	case 32:
		_num := uint32(num)
		return []byte{byte(_num >> 24), byte(_num >> 16), byte(_num >> 8), byte(_num)}
	default:
		return []byte{}
	}
}
