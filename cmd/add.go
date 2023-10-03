package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/magdyamr542/browser-tab-groups/configManager"
	"github.com/urfave/cli/v2"
)

var AddCmd cli.Command = cli.Command{
	Name:        "add",
	Usage:       "Add a url to a tab group",
	Description: "Add a url to a tab group. The tab group will be created if it doesn't exist",
	UsageText: `browser-tab-groups add <tap group path...> <url>

1. Add "https://wwww.google.com" to tab group one two three:
		browser-tab-groups add one two three https://wwww.google.com
`,
	Action: func(cCtx *cli.Context) error {
		jsonCmg, err := configManager.NewJsonConfigManager()
		if err != nil {
			return err
		}

		groupsThenUrl := cCtx.Args().Slice()
		if len(groupsThenUrl) == 0 {
			return fmt.Errorf("you need to provide the tab group path then the url")
		}

		if len(groupsThenUrl) == 1 {
			return fmt.Errorf("you need to provide a url to add to the given tab group")
		}

		url := groupsThenUrl[len(groupsThenUrl)-1]
		groups := groupsThenUrl[:len(groupsThenUrl)-1]
		return addUrlToTapGroup(os.Stdout, jsonCmg, url, groups...)

	},
}

func addUrlToTapGroup(outputW io.Writer, cm configManager.ConfigManager, url string, tapGroups ...string) error {
	err := cm.AddUrl(url, tapGroups...)
	if err != nil {
		return err
	}
	outputW.Write([]byte("url added"))
	return nil
}
