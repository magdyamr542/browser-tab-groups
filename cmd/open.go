package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/magdyamr542/browser-tab-groups/browser"
	"github.com/magdyamr542/browser-tab-groups/configManager"
	"github.com/urfave/cli/v2"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

var OpenCmd cli.Command = cli.Command{
	Name:        "open",
	Usage:       "Open a tap group in the browser",
	Description: "Open a tap group in the browser. This opens all urls in the tap group",
	Aliases:     []string{"o", "op"},
	Action: func(cCtx *cli.Context) error {
		jsonCmg, err := configManager.NewJsonConfigManager()
		if err != nil {
			return err
		}

		tapGroups := cCtx.Args().Slice()
		if len(tapGroups) == 0 {
			return fmt.Errorf("provide a path to the tap group you want to open (as space separated string)")
		}

		return openTapGroup(os.Stdout, tapGroups, jsonCmg, browser.NewBrowser())
	},
}

func openTapGroup(outputW io.Writer, tapGroups []string, cm configManager.ConfigManager, br browser.Browser) error {

	matcher := func(tapGroupPath []string) bool {
		if len(tapGroups) != len(tapGroupPath) {
			return false
		}

		for i, tapGroupLike := range tapGroups {
			savedTapGroup := tapGroupPath[i]
			if !fuzzy.Match(strings.ToLower(tapGroupLike), savedTapGroup) {
				return false
			}
		}

		return true
	}

	tgs, err := cm.GetMatchingTapGroups(matcher)
	if err != nil {
		return err
	}

	urls := make([]string, 0)

	for i := range tgs {
		subUrls, err := tgs[i].Urls()
		if err != nil {
			return err
		}
		urls = append(urls, subUrls...)
	}

	if len(urls) == 0 {
		return fmt.Errorf("no matching for the given pattern %s", tapGroups)
	}

	return br.OpenLinks(urls)
}
