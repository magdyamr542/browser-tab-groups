package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/magdyamr542/browser-tab-groups/configManager"
)

// Listing tap groups with their urls
type LsCmd struct {
	OnlyUrls  bool     `help:"output urls without tap groups"`
	TapGroups []string `arg:"" optional:"" name:"path to tap group" help:"the path to the tap group to remove"`
}

func (ls *LsCmd) Run() error {

	jsonCmg, err := configManager.NewJsonConfigManager()
	if err != nil {
		return err
	}
	return listTapGroups(os.Stdout, jsonCmg, ls.TapGroups, ls.OnlyUrls)
}

// listTapGroups lists all tap groups
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
