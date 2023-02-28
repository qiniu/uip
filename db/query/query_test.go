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
	ip := net.ParseIP("220.248.53.1")
	info, mask, err := db.Query(ip)
	assert.Nil(t, err)
	assert.Less(t, 16, mask)
	assert.NotNil(t, info)
	log.Println(db.Query(ip))
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
	info, mask, err := db.Query(ip)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Greater(t, 16, mask)
	log.Println(db.Query(ip))
}

// only search ipdb
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
//BenchmarkDb_Query
//BenchmarkDb_Query-8   	 2558755	       52.63 ns/op

func BenchmarkDb_Query(b *testing.B) {
	ipdbPath := os.Getenv("IPDBv4_PATH")
	if ipdbPath == "" {
		return
	}

	db, err := NewDb(ipdbPath)
	if err != nil {
		return
	}
	ip := net.ParseIP("220.248.53.1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err = db.Query(ip)
	}
}

// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkDbv6_Query
// BenchmarkDbv6_Query-8   	 2172760	       319.8 ns/op
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
		_, _, err = db.Query(ip)
	}
}
