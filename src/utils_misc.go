package main

import (
	"log"
	"os"

	"encoding/json"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

func InSlice[T comparable](slice []T, elem interface{}) bool {
	for _, v := range slice {
		if elem == v {
			return true
		}
	}
	return false
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
