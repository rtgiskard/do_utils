package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"github.com/pelletier/go-toml"
)

type C_ZT struct {
	Token   string
	UID     string        `yaml:",omitempty"`
	URL     string        `yaml:",omitempty"`
	Timeout time.Duration `yaml:",omitempty"`
	Net     ZtNetPost
	Netm    ZtNetMemberPost
}

type C_DO struct {
	Token string

	Droplet struct {
		Name     string
		OS       string
		Key      string
		Size     string
		Region   []string
		Tags     []string
		UserData string
	}

	ListOption godo.ListOptions
}

type Config struct {
	Zerotier     C_ZT
	DigitalOcean C_DO
}

func (c *Config) load(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = toml.Unmarshal(content, c)
	if err != nil {
		log.Fatal(err)
	}

	// convert to seconds
	c.Zerotier.Timeout *= time.Second
}

func (c *Config) show(format string) {
	fmt.Println(dumps(c, format))
}
