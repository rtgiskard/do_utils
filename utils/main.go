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
		switch args.Zt.Op {
		case "info":
			zt_user_record()
		case "net_ls":
			zt_network_list()
		case "net_add":
			zt_network_create()
		case "net_set":
			zt_network_set(args.Zt.Nid)
		case "net_rm":
			zt_network_del(args.Zt.Nid)
		case "netm_ls":
			zt_network_member_list(args.Zt.Nid)
		case "netm_set":
			zt_network_member_set(args.Zt.Nid, args.Zt.Mid)
		case "netm_rm":
			zt_network_member_del(args.Zt.Nid, args.Zt.Mid)
		default:
			fmt.Println("** Undefined operation:", args.Zt.Op)
		}

	case args.Do != nil:
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
}
