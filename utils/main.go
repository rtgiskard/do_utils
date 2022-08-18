package main

import (
	"os"

	"github.com/alexflint/go-arg"
)

var config = utilsConfig{}

func main() {
	p := arg.MustParse(&args)

	// load config after parse but before sync
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
