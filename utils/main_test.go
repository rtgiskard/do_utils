package main

import (
	"os"
	"testing"

	"github.com/alexflint/go-arg"
	"github.com/stretchr/testify/assert"
)

func TestArgParse(t *testing.T) {
	os.Args = []string{"./utils", "info", "--fmt", "yaml"}
	p := arg.MustParse(&args)
	assert.NotNil(t, p)
}

func TestMainCmd(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		os.Args = []string{"", "info", "--fmt", ""}
		main()
	})

	t.Run("info_with_fmt", func(t *testing.T) {
		for _, fmt := range []string{"toml", "yaml", "json", ""} {
			os.Args = []string{"", "info", "--fmt", fmt}
			p := arg.MustParse(&args)
			assert.NotNil(t, p)

			config.load(args.ConfigFile)
			config.show(args.Format)
		}
	})

	t.Run("zerotier", func(t *testing.T) {
		for _, cmd := range []string{"info", "net_ls", "netm_ls"} {
			os.Args = []string{"", "zerotier", cmd, "-v", "--fmt", ""}
			p := arg.MustParse(&args)
			assert.NotNil(t, p)

			config.load(args.ConfigFile)
			syncArgsZt()
			runZerotierCmd()
		}
	})

	t.Run("digitalocean", func(t *testing.T) {
		for _, cmd := range []string{"ls", "add"} {
			os.Args = []string{"", "digitalocean", cmd, "--name", "abcd", "-n", "-v", "--fmt", ""}
			p := arg.MustParse(&args)
			assert.NotNil(t, p)

			config.load(args.ConfigFile)
			syncArgsDo()
			runDigitaloceanCmd()
		}
	})
}

func TestZtDisplay(t *testing.T) {
	// reset the global args
	args.Verbose = false

	t.Run("displayNetworks", func(t *testing.T) {
		s := make([]ZtNetInfo, 4)
		displayNetworks(s)
	})

	t.Run("displayNetworkMembers", func(t *testing.T) {
		s := make([]ZtNetMemberInfo, 4)
		displayNetworkMembers(s)
	})
}
