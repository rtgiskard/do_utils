package main

import (
	"log"
	"os"
)

func InSlice[T comparable](slice []T, elem interface{}) bool {
	for _, v := range slice {
		if elem == v {
			return true
		}
	}
	return false
}

func ReadFile(path string, bytes int) ([]byte, int) {
	file, err := os.Open(userdata)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	buf := make([]byte, 64*1024)
	n, err := file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	return buf, n
}
