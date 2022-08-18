package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

type doCmd struct {
	Op       string `arg:"positional,required" help:"ls|add|rm|reboot|poweron|poweroff|powercycle"`
	Name     string `arg:"--name" help:"name of the droplet to operate"`
	Size     string `arg:"--size" help:"size of the new droplet"`
	Userdata string `arg:"--userdata" default:"none" help:"source of userdata: [none|gen|$file]"`
	Helper   string `arg:"--helper" default:"../tool/01_gen_userdata.sh" help:"helper to generate userdata"`
	Token    string `arg:"--token" help:"set the api token"`
}

type ztCmd struct {
	Op      string `arg:"positional,required" help:"info|net_ls|net_add|net_set|net_rm|netm_ls|netm_set|netm_rm"`
	UID     string `arg:"--uid" help:"user id"`
	NID     string `arg:"--nid" placeholder:"ID" help:"network id"`
	MID     string `arg:"--mid" placeholder:"ID" help:"network member id"`
	Token   string `arg:"--token" help:"specify the api token"`
	Name    string `arg:"--name" help:"name to be set"`
	IP      string `arg:"--ip" help:"set ip for the node"`
	Timeout int    `arg:"--timeout" default:"10" help:"http timeout in seconds"`
}

type infoCmd struct{}

var args struct {
	Do         *doCmd   `arg:"subcommand:digitalocean" help:"operate with digitalocean api"`
	Zt         *ztCmd   `arg:"subcommand:zerotier" help:"operate with zerotier api"`
	Info       *infoCmd `arg:"subcommand:info" help:"dump config info"`
	Noop       bool     `arg:"-n,--dry-run" help:"no action"`
	ConfigFile string   `arg:"-c,--" default:"config.toml" help:"path of config file"`
	Format     string   `arg:"-f,--fmt" default:"toml" help:"output format: yaml|toml|json"`
	Verbose    bool     `arg:"-v,--verbose" defalt:"false" help:"show verbose info"`
}

func runZerotierCmd() {

	c := ZtClient{
		token:   config.Zerotier.Token,
		baseURL: config.Zerotier.URL,
		timeout: config.Zerotier.Timeout,
		fmt:     args.Format,
	}

	c.Init()

	if args.Noop {
		switch args.Zt.Op {
		case "net_add", "net_set":
			fmt.Printf(">> post args:\n%s\n", Dumps(config.Zerotier.Net, args.Format))
			return
		case "netm_set":
			fmt.Printf(">> post args:\n%s\n", Dumps(config.Zerotier.Netm, args.Format))
			return
		}
	}

	switch args.Zt.Op {
	case "info":
		c.DumpUserRecord()
	case "net_ls":
		c.ListNetwork()
	case "net_add":
		c.CreateNetwork(config.Zerotier.Net)
	case "net_set":
		c.SetNetwork(args.Zt.NID, config.Zerotier.Net)
	case "net_rm":
		c.DelNetwork(args.Zt.NID)
	case "netm_ls":
		c.ListNetworkMember(args.Zt.NID)
	case "netm_set":
		c.SetNetworkMember(args.Zt.NID, args.Zt.MID, config.Zerotier.Netm)
	case "netm_rm":
		c.DelNetworkMember(args.Zt.NID, args.Zt.MID)
	default:
		fmt.Println("** Undefined operation:", args.Zt.Op)
	}
}

func runDigitaloceanCmd() {
	c := DoClient{
		args: &config.DigitalOcean,
	}

	c.Init()

	switch args.Do.Op {
	case "ls":
		c.ListDroplet()
	case "add":
		c.CreateDroplet(args.Noop)
	case "rm":
		c.DestroyDroplet(c.args.Droplet.Name)
	case "reboot", "poweron", "poweroff", "powercycle":
		c.DropletAction(args.Do.Op)
	default:
		fmt.Println("** Undefined operation:", args.Do.Op)
	}
}

func syncArgsZt() bool {
	if args.Zt.Token != "" {
		config.Zerotier.Token = args.Zt.Token
	}
	if args.Zt.UID != "" {
		config.Zerotier.UID = args.Zt.UID
	}
	if args.Zt.IP != "" {
		config.Zerotier.Netm.Config.IPAssignments = []string{args.Zt.IP}
	}
	if args.Zt.Timeout > 0 {
		config.Zerotier.Timeout = time.Duration(args.Zt.Timeout) * time.Second
	}

	switch args.Zt.Op {
	case "net_set", "net_rm":
		if args.Zt.NID == "" {
			fmt.Println("** network id not specified!")
			return false
		}
	case "netm_set", "netm_rm":
		if args.Zt.NID == "" || args.Zt.MID == "" {
			fmt.Println("** network id or member id not specified!")
			return false
		}
	}

	if args.Zt.Name != "" {
		switch args.Zt.Op {
		case "net_set", "net_add":
			config.Zerotier.Net.Config.Name = args.Zt.Name
		case "netm_set", "netm_rm":
			config.Zerotier.Netm.Name = args.Zt.Name
		}
	}

	return true
}

func syncArgsDo() bool {
	if args.Do.Token != "" {
		config.DigitalOcean.Token = args.Do.Token
	}

	if args.Do.Name != "" {
		config.DigitalOcean.Droplet.Name = args.Do.Name
		config.DigitalOcean.Droplet.Tags = []string{args.Do.Name}
	}

	if args.Do.Size != "" {
		config.DigitalOcean.Droplet.Size = args.Do.Size
	}

	switch args.Do.Userdata {
	case "none": // disable userdata
		config.DigitalOcean.Droplet.UserData = ""
	case "gen": // generate userdata
		out, err := exec.Command(args.Do.Helper).Output()
		if err != nil {
			log.Fatal(err)
		}
		config.DigitalOcean.Droplet.UserData = string(out)
	default: // userdata is a file, read up to 64k
		buf, n := ReadFile(args.Do.Userdata, 64*1024)
		config.DigitalOcean.Droplet.UserData = string(buf[:n])
	}

	return true
}
