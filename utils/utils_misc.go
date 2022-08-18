package main

import (
	"log"
	"math/rand"
	"os"
	"strings"

	"encoding/json"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

func InSlice[T comparable](s []T, e interface{}) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func GetStrSet(n int) string {
	strSet := []string{
		"0123456789",
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
	}

	if n < 0 {
		n = 0
	} else if n > len(strSet) {
		n = len(strSet)
	}

	return strings.Join(strSet[:n], "")
}

func ReprBitsLen(num uint64) int {
	for i := 1; ; i++ {
		num >>= 1

		if num == 0 {
			return i
		}
	}
}

func GenRandStr[T string | int](n int, s T) string {

	// reseed should be performed
	// rand.Seed(time.Now().UnixNano())

	var ss string

	// get source set of characters
	var i interface{} = s
	if v, ok := i.(string); ok {
		ss = v
	} else if v, ok := i.(int); ok {
		ss = GetStrSet(v)
	}

	// check length
	if len(ss) == 0 || n == 0 {
		return ""
	}

	// rand select
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(ss[rand.Intn(len(ss))])
	}
	return sb.String()
}

func ReadFile(path string, size int) ([]byte, int) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	buf := make([]byte, size)
	n, err := file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	return buf, n
}

func dumps(o interface{}, format string) string {
	var b []byte
	var err error

	switch format {
	case "yaml":
		b, err = yaml.Marshal(o)
	case "json":
		b, err = json.MarshalIndent(o, "", "\t")
	case "toml":
		b, err = toml.Marshal(o)
	default:
		return ""
	}

	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}
