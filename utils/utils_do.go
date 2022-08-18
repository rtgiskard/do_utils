package main

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/digitalocean/godo"
)

type DoClient struct {
	client *godo.Client
	ctx    context.Context
	args   *cDO
}

func (c *DoClient) Init() {
	c.client = godo.NewFromToken(c.args.Token)
	c.ctx = context.TODO()
}

func (c *DoClient) GetSSHKeyFp(keyname string) string {
	keys, _, err := c.client.Keys.List(c.ctx, &c.args.ListOption)

	if err != nil {
		log.Fatal(err)
	}

	for _, key := range keys {
		if key.Name == keyname {
			return key.Fingerprint
		}
	}

	return ""
}

func (c *DoClient) GetDropletRegion() string {
	regions, _, _ := c.client.Regions.List(c.ctx, &c.args.ListOption)

	for _, region := range regions {
		if !region.Available {
			continue
		}

		if InSlice(c.args.Droplet.Region, region.Slug) &&
			InSlice(region.Sizes, c.args.Droplet.Size) {
			return region.Slug
		}
	}

	return ""
}

func (c *DoClient) GetDropletImage() string {
	images, _, _ := c.client.Images.List(c.ctx, &c.args.ListOption)

	slugs := make([]string, 0, 2)

	for _, image := range images {
		if image.Distribution == c.args.Droplet.OS {
			slugs = append(slugs, image.Slug)
		}
	}

	// select the lastest version
	sort.Strings(slugs)

	if len(slugs) > 0 {
		return slugs[len(slugs)-1]
	}

	return ""
}

func (c *DoClient) DropletAction(action string) {
	id := c.GetDropletID(c.args.Droplet.Name)
	if id == 0 {
		log.Fatal("failed to get droplet id")
		return
	}

	actionMap := map[string]func(context.Context, int) (*godo.Action, *godo.Response, error){
		"reboot":     c.client.DropletActions.Reboot,
		"poweron":    c.client.DropletActions.PowerOn,
		"poweroff":   c.client.DropletActions.PowerOff,
		"powercycle": c.client.DropletActions.PowerCycle,
	}

	fmt.Printf("-> droplet action: %s %s\n", action, c.args.Droplet.Name)
	actionMap[action](c.ctx, id)
}

func (c *DoClient) GetDropletID(name string) int {
	droplets, _, _ := c.client.Droplets.List(c.ctx, &c.args.ListOption)

	for _, droplet := range droplets {
		if droplet.Name == name {
			return droplet.ID
		}
	}
	return 0
}

func (c *DoClient) ListDroplet() {
	droplets, _, _ := c.client.Droplets.List(c.ctx, &c.args.ListOption)

	fmt.Println("-> list droplet:")
	for _, droplet := range droplets {

		ipv4 := ""
		ipv6 := ""

		for _, net := range droplet.Networks.V4 {
			if net.Type == "public" {
				ipv4 = net.IPAddress
			}
		}
		for _, net := range droplet.Networks.V6 {
			if net.Type == "public" {
				ipv6 = net.IPAddress
			}
		}

		fmt.Printf("%s: %s %s (%s, %s, %s)\n", droplet.Name, ipv4, ipv6,
			droplet.Region.Slug, droplet.SizeSlug, droplet.Status)
	}
}

func (c *DoClient) CreateDroplet(noop bool) {
	chRegion := make(chan string)
	chImage := make(chan string)
	chKey := make(chan string)
	chID := make(chan int)

	go func() {
		chID <- c.GetDropletID(c.args.Droplet.Name)
	}()

	go func() {
		chRegion <- c.GetDropletRegion()
	}()

	go func() {
		chImage <- c.GetDropletImage()
	}()

	go func() {
		chKey <- c.GetSSHKeyFp(c.args.Droplet.Key)
	}()

	// abort if droplet exist
	if <-chID != 0 {
		fmt.Printf("Abort on existing droplet: %s\n", c.args.Droplet.Name)
		return
	}

	// verify: region, image, sshKey

	region := <-chRegion
	if region == "" {
		fmt.Printf("Invalid region: %s\n", c.args.Droplet.Region)
		return
	}

	var image godo.DropletCreateImage
	if slug := <-chImage; slug != "" {
		image = godo.DropletCreateImage{Slug: slug}
	} else {
		fmt.Printf("Invalid OS: %s\n", c.args.Droplet.OS)
		return
	}

	var sshKey []godo.DropletCreateSSHKey
	if sshKeyFp := <-chKey; sshKeyFp != "" {
		sshKey = []godo.DropletCreateSSHKey{{Fingerprint: sshKeyFp}}
	} else {
		fmt.Printf("Invalid ssh key: %s\n", c.args.Droplet.Key)
		return
	}

	createRequest := &godo.DropletCreateRequest{
		Name:     c.args.Droplet.Name,
		Size:     c.args.Droplet.Size,
		Region:   region,
		Image:    image,
		SSHKeys:  sshKey,
		Backups:  false,
		IPv6:     true,
		Tags:     []string{c.args.Droplet.Name},
		UserData: c.args.Droplet.UserData,
	}

	fmt.Printf("-> create droplet:\n--\n%s\n", dumps(createRequest, "toml"))

	if !noop {
		c.client.Droplets.Create(c.ctx, createRequest)
	}
}

func (c *DoClient) DestroyDroplet(name string) {
	fmt.Println("-> delete droplet:", name)
	c.client.Droplets.DeleteByTag(c.ctx, name)
}
