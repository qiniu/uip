package data

import (
	"bufio"
	"github.com/qiniu/uip/db/field"
	"log"
	"strings"
)

var NormalizedMap = make(map[string]map[string]string)

func parseData(data string) map[string]string {
	m := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		split := strings.SplitN(string(line), "\t", 2)
		if len(split) < 2 {
			continue
		}
		m[split[0]] = split[1]
	}
	if err := scanner.Err(); err != nil {
		log.Println("load data failed", err)
		return nil
	}
	return m
}

// Load will overwrite the old data
func Load(key, data string) {
	NormalizedMap[key] = parseData(data)
}

func LoadDefault() {
	Load(field.City, CityList)
	Load(field.Province, ProvinceList)
	Load(field.ISP, ISPList)
	Load(field.Continent, ContinentList)
}

func init() {
	LoadDefault()
}
