package main

import (
	"log"
	"os"

	"github.com/magdyamr542/browser-tab-groups/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "browser-tab-groups",
		Usage:       "Group urls and open them together",
		Description: "Save grouped urls from the command line. Open multiple urls, which are grouped together. Or a single url using its group path",
		Commands: []*cli.Command{
			&cmd.AddCmd,
			&cmd.LsCmd,
			&cmd.EditCmd,
			&cmd.RmCmd,
			&cmd.OpenCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
