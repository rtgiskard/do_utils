package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"github.com/pelletier/go-toml/v2"
)

type cZT struct {
	Token   string
	UID     string        `yaml:",omitempty"`
	URL     string        `yaml:",omitempty"`
	Timeout time.Duration `yaml:",omitempty"`
	Net     ZtNetPost
	Netm    ZtNetMemberPost
}

type cDO struct {
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

type mainConfig struct {
	Zerotier     cZT
	DigitalOcean cDO
}

func (c *mainConfig) load(path string) {
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

func (c *mainConfig) show(format string) {
	fmt.Println(Dumps(c, format))
}
