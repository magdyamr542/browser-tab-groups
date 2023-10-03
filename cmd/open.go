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
	Usage:       "Open a tab group in the browser",
	Description: "Open a tab group in the browser. This opens all urls in the tab group",
	Aliases:     []string{"o", "op"},
	Action: func(cCtx *cli.Context) error {
		jsonCmg, err := configManager.NewJsonConfigManager()
		if err != nil {
			return err
		}

		tabGroups := cCtx.Args().Slice()
		if len(tabGroups) == 0 {
			return fmt.Errorf("provide a path to the tab group you want to open (as space separated string)")
		}

		return openTabGroup(os.Stdout, tabGroups, jsonCmg, browser.NewBrowser())
	},
}

func openTabGroup(outputW io.Writer, tabGroups []string, cm configManager.ConfigManager, br browser.Browser) error {

	matcher := func(tgPath []string) bool {
		if len(tabGroups) != len(tgPath) {
			return false
		}

		for i, tgLike := range tabGroups {
			savedTabGroup := tgPath[i]
			if !fuzzy.Match(strings.ToLower(tgLike), savedTabGroup) {
				return false
			}
		}

		return true
	}

	tgs, err := cm.GetMatchingTabGroups(matcher)
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
		return fmt.Errorf("no matching for the given pattern %s", tabGroups)
	}

	return br.OpenLinks(urls)
}
