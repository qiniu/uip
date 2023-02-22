package query

import (
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"os"
	"testing"

	_ "github.com/qiniu/uip/db/format/ipdb"
)

func TestDb_Query(t *testing.T) {
	ipdbPath := os.Getenv("IPDBv4_PATH")
	if ipdbPath == "" {
		return
	}
	db, err := NewDb(ipdbPath)
	if err != nil {
		return
	}
	ip := net.ParseIP("220.248.53.61")
	info, err := db.Query(ip)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	log.Println(db.Query(ip))
	assert.Equal(t, "中国", info.Country)
	assert.NotNil(t, db.VersionInfo())
	assert.True(t, db.VersionInfo().HasIpV4())
}

func TestDbv6_Query(t *testing.T) {
	ipdbPath := os.Getenv("IPDBv6_PATH")
	if ipdbPath == "" {
		return
	}
	db, err := NewDb(ipdbPath)
	if err != nil {
		return
	}
	ip := net.ParseIP("2001:4860:4860::8888")
	info, err := db.Query(ip)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	log.Println(db.Query(ip))
	assert.Equal(t, "美国", info.Country)
	assert.NotNil(t, db.VersionInfo())
	assert.True(t, db.VersionInfo().HasIpV6())
}

// only search ipdb
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// 	BenchmarkDb_Query
//	BenchmarkDb_Query-8   	 4014349	       289.3 ns/op
// split search result to fields
//BenchmarkDb_Query
//BenchmarkDb_Query-8   	 2558755	       461.4 ns/op

func BenchmarkDb_Query(b *testing.B) {
	ipdbPath := os.Getenv("IPDBv4_PATH")
	if ipdbPath == "" {
		return
	}
	ipdbPath = "/Users/long/github/qiniu/uip/ipli.ipdb"
	db, err := NewDb(ipdbPath)
	if err != nil {
		return
	}
	ip := net.ParseIP("220.248.53.61")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = db.Query(ip)
	}
}

// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkDbv6_Query
// BenchmarkDbv6_Query-8   	 2172760	       528.3 ns/op
func BenchmarkDbv6_Query(b *testing.B) {
	ipdbPath := os.Getenv("IPDBv6_PATH")
	if ipdbPath == "" {
		return
	}

	db, err := NewDb(ipdbPath)
	if err != nil {
		return
	}
	ip := net.ParseIP("2001:4860:4860::8888")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = db.Query(ip)
	}
}

