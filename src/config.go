package main

import (
	"io/ioutil"
	"log"

	"github.com/digitalocean/godo"
	"github.com/pelletier/go-toml"
)

type conf_auth struct {
	Token string
}

type conf_droplet struct {
	Name     string
	OS       string
	Key      string
	Size     string
	Region   []string
	Tags     []string
	UserData string
}

type Config struct {
	noop       bool
	Auth       conf_auth
	Droplet    conf_droplet
	ListOption godo.ListOptions
}

func (c *Config) load_config(path string) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	err = toml.Unmarshal(content, c)

	if err != nil {
		log.Fatal(err)
	}
}
