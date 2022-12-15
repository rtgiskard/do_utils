package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/digitalocean/godo"
)

// DoClient is a wrapper of general operation for digitalocean
type DoClient struct {
	client *godo.Client
	ctx    context.Context
	args   *cDO
}

// Init setup the backend client and context
func (c *DoClient) Init() {
	c.client = godo.NewFromToken(c.args.Token)
	c.ctx = context.TODO()
}

// GetSSHKeyFp returns fingerprint for the sshkey
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

// GetDropletRegion returns the prefered available region which match the args
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

// GetDropletImage returns the slug name for the distro of latest version
func (c *DoClient) GetDropletImage(distro string) string {
	images, _, _ := c.client.Images.List(c.ctx, &c.args.ListOption)

	slugs := make([]string, 0, 2)

	for _, image := range images {
		if image.Slug == distro {
			return distro
		}
		if strings.ToLower(image.Distribution) == strings.ToLower(distro) {
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

// DropletAction perform the basic droplet action like reboot and poweroff
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

// GetDropletID returns id of the named droplet
func (c *DoClient) GetDropletID(name string) int {
	droplets, _, _ := c.client.Droplets.List(c.ctx, &c.args.ListOption)

	for _, droplet := range droplets {
		if droplet.Name == name {
			return droplet.ID
		}
	}
	return 0
}

// ListDroplet list all existing droplet
func (c *DoClient) ListDroplet() {
	droplets, _, _ := c.client.Droplets.List(c.ctx, &c.args.ListOption)

	if len(droplets) == 0 {
		fmt.Println("<empty>")
		return
	}

	info := [][]interface{}{
		{"Name", "Ipv4", "Ipv6", "Region", "Size", "Status", "Created"},
	}

	for _, droplet := range droplets {

		ipv4, _ := droplet.PublicIPv4()
		ipv6, _ := droplet.PublicIPv6()

		info = append(info, []interface{}{
			droplet.Name, ipv4, ipv6, droplet.Region.Slug, droplet.SizeSlug,
			droplet.Status, droplet.Created})
	}

	ShowTable(info)
}

// CreateDroplet create new droplet and apply the settings
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
		chImage <- c.GetDropletImage(c.args.Droplet.OS)
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

	fmt.Printf("-> create droplet:\n--\n%s\n", Dumps(createRequest, "toml"))

	if !noop {
		_, _, err := c.client.Droplets.Create(c.ctx, createRequest)

		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}
}

// DestroyDroplet destroy droplet by the pre set tag which equals the name
func (c *DoClient) DestroyDroplet(name string) {
	resp, _ := c.client.Droplets.DeleteByTag(c.ctx, name)
	if resp.StatusCode == 204 {
		fmt.Println("droplet removed:", name)
	} else {
		fmt.Println("Status:", resp.Status)
	}
}
