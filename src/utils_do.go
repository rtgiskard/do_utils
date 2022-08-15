package main

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/digitalocean/godo"
)

func do_get_client(token string) (context.Context, *godo.Client) {
	return context.TODO(), godo.NewFromToken(token)
}

func do_ssh_key_get_fp(ctx context.Context, client *godo.Client, keyname string) string {
	keys, _, err := client.Keys.List(ctx, &config.DigitalOcean.ListOption)

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

func do_datacenter_region_get(ctx context.Context, client *godo.Client) string {
	regions, _, _ := client.Regions.List(ctx, &config.DigitalOcean.ListOption)

	for _, region := range regions {
		if !region.Available {
			continue
		}

		if InSlice(config.DigitalOcean.Droplet.Region, region.Slug) && InSlice(region.Sizes, config.DigitalOcean.Droplet.Size) {
			return region.Slug
		}
	}

	return ""
}

func do_droplet_image_get(ctx context.Context, client *godo.Client) string {
	images, _, _ := client.Images.List(ctx, &config.DigitalOcean.ListOption)

	images_slice := make([]string, 0, 2)

	for _, image := range images {
		if image.Distribution == config.DigitalOcean.Droplet.OS {
			images_slice = append(images_slice, image.Slug)
		}
	}

	// select the lastest version
	sort.Strings(images_slice)

	last_index := len(images_slice) - 1
	if last_index >= 0 {
		return images_slice[last_index]
	}

	return ""
}

func do_droplet_action(ctx context.Context, client *godo.Client, action string) {
	id := do_droplet_get_id(ctx, client, config.DigitalOcean.Droplet.Name)
	if id == 0 {
		log.Fatal("failed to get droplet id")
		return
	}

	action_map := map[string]func(context.Context, int) (*godo.Action, *godo.Response, error){
		"reboot":     client.DropletActions.Reboot,
		"poweron":    client.DropletActions.PowerOn,
		"poweroff":   client.DropletActions.PowerOff,
		"powercycle": client.DropletActions.PowerCycle,
	}

	fmt.Printf("-> droplet action: %s %s\n", action, config.DigitalOcean.Droplet.Name)
	action_map[action](ctx, id)
}

func do_droplet_get_id(ctx context.Context, client *godo.Client, name string) int {
	droplets, _, _ := client.Droplets.List(ctx, &config.DigitalOcean.ListOption)

	for _, droplet := range droplets {
		if droplet.Name == name {
			return droplet.ID
		}
	}
	return 0
}

func do_droplet_ls(ctx context.Context, client *godo.Client) {
	droplets, _, _ := client.Droplets.List(ctx, &config.DigitalOcean.ListOption)

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

func do_droplet_create(ctx context.Context, client *godo.Client) {
	ch_region := make(chan string)
	ch_image := make(chan string)
	ch_key := make(chan string)
	ch_id := make(chan int)

	go func() {
		ch_id <- do_droplet_get_id(ctx, client, config.DigitalOcean.Droplet.Name)
	}()

	go func() {
		ch_region <- do_datacenter_region_get(ctx, client)
	}()

	go func() {
		ch_image <- do_droplet_image_get(ctx, client)
	}()

	go func() {
		ch_key <- do_ssh_key_get_fp(ctx, client, config.DigitalOcean.Droplet.Key)
	}()

	// abort if droplet exist
	if <-ch_id != 0 {
		fmt.Printf("Abort on existing droplet: %s\n", config.DigitalOcean.Droplet.Name)
		return
	}

	// verify: region, image, ssh_key

	region := <-ch_region
	if region == "" {
		fmt.Printf("Invalid region: %s\n", config.DigitalOcean.Droplet.Region)
		return
	}

	var image godo.DropletCreateImage
	if image_slug := <-ch_image; image_slug != "" {
		image = godo.DropletCreateImage{Slug: image_slug}
	} else {
		fmt.Printf("Invalid OS: %s\n", config.DigitalOcean.Droplet.OS)
		return
	}

	var ssh_key []godo.DropletCreateSSHKey
	if ssh_key_fp := <-ch_key; ssh_key_fp != "" {
		ssh_key = []godo.DropletCreateSSHKey{{Fingerprint: ssh_key_fp}}
	} else {
		fmt.Printf("Invalid ssh key: %s\n", config.DigitalOcean.Droplet.Key)
		return
	}

	createRequest := &godo.DropletCreateRequest{
		Name:     config.DigitalOcean.Droplet.Name,
		Size:     config.DigitalOcean.Droplet.Size,
		Region:   region,
		Image:    image,
		SSHKeys:  ssh_key,
		Backups:  false,
		IPv6:     true,
		Tags:     []string{config.DigitalOcean.Droplet.Name},
		UserData: config.DigitalOcean.Droplet.UserData,
	}

	fmt.Printf("-> create droplet:\n--\n%s\n", dumps(createRequest, "toml"))

	if !args.Noop {
		client.Droplets.Create(ctx, createRequest)
	}
}

func do_droplet_rm(ctx context.Context, client *godo.Client) {
	fmt.Println("-> delete droplet:", config.DigitalOcean.Droplet.Name)
	client.Droplets.DeleteByTag(ctx, config.DigitalOcean.Droplet.Name)
}
