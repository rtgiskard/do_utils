package main

import (
	"fmt"
)

var config = Config{}
var config_file = "config.toml"
var userdata_generator = "../tool/01_gen_userdata.sh"

func main() {
	parse_args()

	switch {
	case args.Info != nil:
		config.show(args.Format)

	case args.Zt != nil:
		run_zerotier_cmd()

	case args.Do != nil:
		run_digitalocean_cmd()
	}
}

func run_zerotier_cmd() {

	client := ZtClient{
		token:   config.Zerotier.Token,
		baseURL: config.Zerotier.URL,
		timeout: config.Zerotier.Timeout,
		fmt:     args.Format,
	}

	client.Init()

	if args.Noop {
		switch args.Zt.Op {
		case "net_add", "net_set":
			fmt.Printf(">> post args:\n%s\n", dumps(config.Zerotier.Net, args.Format))
			return
		case "netm_set":
			fmt.Printf(">> post args:\n%s\n", dumps(config.Zerotier.Netm, args.Format))
			return
		}
	}

	switch args.Zt.Op {
	case "info":
		client.DumpUserRecord()
	case "net_ls":
		client.ListNetwork()
	case "net_add":
		client.CreateNetwork(config.Zerotier.Net)
	case "net_set":
		client.SetNetwork(args.Zt.Nid, config.Zerotier.Net)
	case "net_rm":
		client.DelNetwork(args.Zt.Nid)
	case "netm_ls":
		client.ListNetworkMember(args.Zt.Nid)
	case "netm_set":
		client.SetNetworkMember(args.Zt.Nid, args.Zt.Mid, config.Zerotier.Netm)
	case "netm_rm":
		client.DelNetworkMember(args.Zt.Nid, args.Zt.Mid)
	default:
		fmt.Println("** Undefined operation:", args.Zt.Op)
	}
}

func run_digitalocean_cmd() {
	ctx, client := do_get_client(config.DigitalOcean.Token)

	switch args.Do.Op {
	case "ls":
		do_droplet_ls(ctx, client)
	case "add":
		do_droplet_create(ctx, client)
	case "rm":
		do_droplet_rm(ctx, client)
	case "reboot", "poweron", "poweroff", "powercycle":
		do_droplet_action(ctx, client, args.Do.Op)
	default:
		fmt.Println("** Undefined operation:", args.Do.Op)
	}
}
