package main

import (
	"os"

	"github.com/magdyamr542/browser-tab-groups/cmd"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Ls   cmd.LsCmd   `cmd:"" help:"List all tab groups" aliases:"l,ls,list"`
	Add  cmd.AddCmd  `cmd:"" help:"Add a new url to a tap group"`
	Open cmd.OpenCmd `cmd:"" help:"Open a tap group in the browser"`
	Rm   cmd.RmCmd   `cmd:"" help:"Remove a saved tap group"`
}

func main() {
	// If running without any extra arguments, default to the --help flag
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
