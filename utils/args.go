package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/alexflint/go-arg"
)

type DoCmd struct {
	Op       string `arg:"positional,required" help:"ls|add|rm|reboot|poweron|poweroff|powercycle"`
	Name     string `arg:"--name" help:"name of the droplet to operate"`
	Size     string `arg:"--size" help:"size of the new droplet"`
	Userdata string `arg:"--userdata" default:"none" help:"create with userdata: [none|gen|$file]"`
	Token    string `arg:"--token" help:"set the api token"`
}

type ZtCmd struct {
	Op      string `arg:"positional,required" help:"info|net_ls|net_add|net_set|net_rm|netm_ls|netm_set|netm_rm"`
	Uid     string `arg:"--uid" help:"user id"`
	Nid     string `arg:"--nid" placeholder:"ID" help:"network id"`
	Mid     string `arg:"--mid" placeholder:"ID" help:"network member id"`
	Token   string `arg:"--token" help:"specify the api token"`
	Name    string `arg:"--name" help:"name to be set"`
	Ip      string `arg:"--ip" help:"set ip for the node"`
	Timeout int    `arg:"--timeout" default:"10" help:"http timeout in seconds"`
}

type InfoCmd struct{}

var args struct {
	Do      *DoCmd   `arg:"subcommand:do" help:"digitalocean utils"`
	Zt      *ZtCmd   `arg:"subcommand:zt" help:"zerotier utils"`
	Info    *InfoCmd `arg:"subcommand:info" help:"show config info"`
	Noop    bool     `arg:"-n,--dry-run" help:"no action"`
	Format  string   `arg:"-f,--fmt" default:"toml" help:"output format: yaml|toml|json"`
	Verbose bool     `arg:"-v,--verbose" defalt:"false" help:"show verbose info"`
}

func parse_args() {
	p := arg.MustParse(&args)
	sync_ret := true

	config.load(config_file)

	switch {
	case args.Info != nil:
	case args.Zt != nil:
		sync_ret = sync_args_zt()
	case args.Do != nil:
		sync_ret = sync_args_do()
	default:
		p.WriteHelp(os.Stdout)
	}

	if !sync_ret {
		os.Exit(1)
	}
}

func sync_args_zt() bool {
	if args.Zt.Token != "" {
		config.Zerotier.Token = args.Zt.Token
	}
	if args.Zt.Uid != "" {
		config.Zerotier.UID = args.Zt.Uid
	}
	if args.Zt.Ip != "" {
		config.Zerotier.Netm.Config.IPAssignments = []string{args.Zt.Ip}
	}

	switch args.Zt.Op {
	case "net_set", "net_rm":
		if args.Zt.Nid == "" {
			fmt.Println("** network id not specified!")
			return false
		}
	case "netm_set", "netm_rm":
		if args.Zt.Nid == "" || args.Zt.Mid == "" {
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

func sync_args_do() bool {
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
		out, err := exec.Command(userdata_generator).Output()
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

func parse_args_with_flag() {
	// note only

	var operation string
	var name string
	var size string
	var userdata string
	var noop bool

	flag.StringVar(&operation, "op", "", "valid: [ls|add|rm|reboot|poweron|poweroff|powercycle]")
	flag.StringVar(&name, "name", "", "name of the droplet")
	flag.StringVar(&size, "size", "", "size of the droplet")
	flag.StringVar(&userdata, "userdata", "", "create droplet with userdata: [''|auto|file_to_read]")
	flag.BoolVar(&noop, "noop", false, "no real operation for droplet creation")

	flag.Parse()

	if operation == "" {
		CommandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
}
