package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/magdyamr542/browser-tab-groups/colors"
	"github.com/magdyamr542/browser-tab-groups/configManager"
	"github.com/magdyamr542/browser-tab-groups/opener"
	op "github.com/magdyamr542/browser-tab-groups/opener"
	"github.com/urfave/cli/v2"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Opener string

const (
	Shell   Opener = "shell"
	Browser Opener = "browser"
)

const (
	OpenerFlag = "opener"
)

var OpenCmd cli.Command = cli.Command{
	Name:        "open",
	Usage:       "Open a tab group in the browser or shell",
	Description: "Open a tab group in the browser or shell. This opens all urls in the tab group",
	Aliases:     []string{"o", "op"},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  OpenerFlag,
			Value: string(Browser),
			Usage: "The opener to use. One of 'browser' or 'shell'. The 'shell' opener makes HTTP GET requests to the url.",
			Action: func(ctx *cli.Context, s string) error {
				if s != string(Shell) && s != string(Browser) {
					return fmt.Errorf("opener must only be 'browser' or 'shell'")
				}
				return nil
			},
		},
	},
	Action: func(cCtx *cli.Context) error {
		jsonCmg, err := configManager.NewJsonConfigManager()
		if err != nil {
			return err
		}

		tabGroups := cCtx.Args().Slice()
		if len(tabGroups) == 0 {
			return fmt.Errorf("provide a path to the tab group you want to open (as space separated string)")
		}

		var linkOpener op.Opener
		opener := Opener(cCtx.String(OpenerFlag))
		if opener == Browser {
			linkOpener = op.NewBrowser()
		} else if opener == Shell {
			linkOpener = op.NewShell(os.Stdout)
		} else {
			return fmt.Errorf("invalid opener flag. this should never happen")
		}

		return openTabGroup(os.Stdout, tabGroups, jsonCmg, linkOpener)
	},
}

func openTabGroup(outputW io.Writer, tabGroups []string, cm configManager.ConfigManager, linkOpener opener.Opener) error {

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

	requests := make([]opener.Request, 0)

	for i := range tgs {
		tgPath := tgs[i].Path()
		tgPathStr := strings.Join(tgPath, " -> ")

		subUrls, err := tgs[i].Urls()
		if err != nil {
			return err
		}
		for _, url := range subUrls {
			desc := tgPathStr + " -> " + url
			requests = append(requests, op.Request{
				Link:        url,
				Description: fmt.Sprintf("%s\n", colors.Bold(desc)),
			})
		}
	}

	if len(requests) == 0 {
		return fmt.Errorf("no matching for the given pattern %s", tabGroups)
	}

	return linkOpener.OpenLinks(requests)
}
