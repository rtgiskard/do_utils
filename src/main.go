package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var config = Config{}
var config_file = "config.toml"
var userdata_generator = "../tool/01_gen_userdata.sh"

var operation string
var name string
var size string
var userdata string
var noop bool

func init() {

	flag.StringVar(&operation, "op", "", "valid: [ls|add|rm|reboot|poweron|poweroff|powercycle]")
	flag.StringVar(&name, "name", "", "name of the droplet")
	flag.StringVar(&size, "size", "", "size of the droplet")
	flag.StringVar(&userdata, "userdata", "", "create droplet with userdata: [''|auto|file_to_read]")
	flag.BoolVar(&noop, "noop", false, "no real operation for droplet creation")
}

func handle_args() {
	flag.Parse()

	config.load_config(config_file)

	if name != "" {
		config.Droplet.Name = name
		config.Droplet.Tags = []string{name}
	}

	if size != "" {
		config.Droplet.Size = size
	}

	switch userdata {
	case "":
		// disable userdata
		config.Droplet.UserData = ""
	case "auto":
		// generate userdata
		out, err := exec.Command(userdata_generator).Output()
		if err != nil {
			log.Fatal(err)
		}
		config.Droplet.UserData = string(out)
	default:
		// userdata is a file, read up to 64k
		buf, n := ReadFile(userdata, 64*1024)
		config.Droplet.UserData = string(buf[:n])
	}

	config.noop = noop

	if operation == "" {
		CommandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	handle_args()

	ctx, client := do_get_client(config.Auth.Token)

	switch operation {
	case "ls":
		do_droplet_ls(ctx, client)
	case "add":
		do_droplet_create(ctx, client)
	case "rm":
		do_droplet_rm(ctx, client)
	case "reboot":
		fallthrough
	case "poweron":
		fallthrough
	case "poweroff":
		fallthrough
	case "powercycle":
		do_droplet_action(ctx, client, operation)
	default:
		fmt.Println("Undefined operation:", operation)
	}
}
