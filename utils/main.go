/*
This is a set of tools for several platforms, which makes it possible to
perform some basic operations without logging into each platform with a browser
*/
package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"
)

var args = mainArgs{}
var config = mainConfig{}

func load_config(path string) bool {
	if IsFileExist(path) {
		config.load(path)
		return true
	} else {
		conf_list := []string{"config.toml", "conf/config.toml"}
		for _, cpath := range conf_list {
			if cpath != path && IsFileExist(cpath) {
				// log.Printf("using fallback config: %s", cpath)
				config.load(cpath)
				return true
			}
		}
	}

	return false
}

func main() {
	p := arg.MustParse(&args)

	if !load_config(args.ConfigFile) {
		log.Fatal("failed to load config file, abort!")
	}

	switch {
	case args.Info != nil:
		config.show(args.Format)

	case args.Zt != nil:
		if !syncArgsZt() {
			os.Exit(1)
		}
		runZerotierCmd()

	case args.Do != nil:
		if !syncArgsDo() {
			os.Exit(1)
		}
		runDigitaloceanCmd()

	default:
		p.WriteHelp(os.Stdout)
	}
}
