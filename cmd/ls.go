package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/magdyamr542/browser-tab-groups/configManager"
	"github.com/urfave/cli/v2"
)

const (
	OnlyUrlsFlag = "only-urls"
)

var LsCmd cli.Command = cli.Command{
	Name:        "list",
	Aliases:     []string{"l", "ls"},
	Usage:       "List all tab groups",
	Description: "List all saved tab groups or provide a path and list only the tap groups which fuzzy match the provided path.",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  OnlyUrlsFlag,
			Value: false,
			Usage: "Show only the urls without the tap group hierarchy",
		},
	},
	UsageText: `browser-tab-groups list [command options] [tap group path...]

1. List all tab groups:
		browser-tab-groups list

2. List tab groups with using a path with fuzzy matching:
		browser-tab-groups list frst scnd thrd
`,
	Action: func(cCtx *cli.Context) error {

		jsonCmg, err := configManager.NewJsonConfigManager()
		if err != nil {
			return err
		}

		// The user can filter for certain entries using Fuzzy matching.
		tapGroups := cCtx.Args().Slice()
		onlyUrls := cCtx.Bool(OnlyUrlsFlag)

		return listTapGroups(os.Stdout, jsonCmg, tapGroups, onlyUrls)
	},
}

func listTapGroups(outputW io.Writer, cm configManager.ConfigManager, tapGroups []string, onlyUrls bool) error {
	matchAll := false
	if len(tapGroups) == 0 {
		matchAll = true
	}

	matcher := func(tapGroupPath []string) bool {
		if matchAll {
			// Match the roots. They cover all children in their urls.
			return len(tapGroupPath) == 1
		}

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

	if len(tgs) == 0 {
		return fmt.Errorf("no matching for the given pattern %s", tapGroups)
	}

	if onlyUrls {
		urls := make([]string, 0)
		for i := range tgs {
			subUrls, err := tgs[i].Urls()
			if err != nil {
				return err
			}
			urls = append(urls, subUrls...)
		}

		for i := range urls {
			outputW.Write([]byte(urls[i] + "\n"))
		}

		return nil
	}

	for i, tg := range tgs {
		str, err := tg.String("")
		if err != nil {
			return err
		}

		outputW.Write([]byte(str + "\n"))
		if i != len(tgs)-1 {
			outputW.Write([]byte("\n"))
		}
	}

	return nil
}
