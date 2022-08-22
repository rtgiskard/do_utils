/*
This is a set of tools for several platforms, which makes it possible to
perform some basic operations without logging into each platform with a browser
*/
package main

import (
	"os"

	"github.com/alexflint/go-arg"
)

var args = mainArgs{}
var config = mainConfig{}

func main() {
	p := arg.MustParse(&args)

	config.load(args.ConfigFile)

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
