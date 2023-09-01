package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/magdyamr542/browser-tab-groups/browser"
	"github.com/magdyamr542/browser-tab-groups/configManager"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// Adding a new tap group
type OpenCmd struct {
	TapGroups []string `arg:"" name:"tap groups" help:"the path to the tap group to open"`
}

func (open *OpenCmd) Run() error {
	jsonCmg, err := configManager.NewJsonConfigManager()
	if err != nil {
		return err
	}
	return openTapGroup(os.Stdout, open.TapGroups, jsonCmg, browser.NewBrowser())
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

	urls, err := cm.GetMatchingUrls(matcher)
	if err != nil {
		return err
	}

	if len(urls) == 0 {
		return fmt.Errorf("no matching for the given pattern %s", tapGroups)
	}

	return br.OpenLinks(urls)
}
