package cmd

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/magdyamr542/browser-tab-groups/browser"
	"github.com/magdyamr542/browser-tab-groups/configManager"
	"github.com/magdyamr542/browser-tab-groups/helpers"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// Adding a new tap group
type OpenCmd struct {
	TapGroup string `arg:"" name:"tap group" help:"the tap group to add the url to"`
	UrlLike  string `arg:"" optional:"" name:"url part" help:"a part of the url to be use with fuzzy matching"`
}

func (open *OpenCmd) Run() error {
	jsonCmg, err := configManager.NewJsonConfigManager()
	if err != nil {
		return err
	}
	return openTapGroup(os.Stdout, open.TapGroup, open.UrlLike, jsonCmg, browser.NewBrowser())
}

// AddUrlToTapGroup adds the given url to the given tap group
func openTapGroup(outputW io.Writer, tapGroup string, urlLike string, cm configManager.ConfigManager, br browser.Browser) error {
	tapGroups, err := cm.GetTapGroups()
	if err != nil {
		return err
	}
	if len(tapGroups) == 0 {
		return errors.New("the given tap group does not exist")
	}

	tapGroupLikeLower := strings.ToLower(tapGroup)
	tapGroups = helpers.Filter(tapGroups, func(tg string) bool {
		return fuzzy.Match(tapGroupLikeLower, strings.ToLower(tg))
	})

	if len(tapGroups) == 0 {
		return errors.New("no matches found in the saved tap groups")
	}

	if len(tapGroups) > 1 {
		return errors.New("more than one tap group matched")
	}

	urls, err := cm.GetUrls(tapGroups[0])
	if err != nil {
		return err
	}
	if len(urls) == 0 {
		return errors.New("the given tap group does not have urls")
	}

	if len(strings.TrimSpace(urlLike)) > 0 {
		urlLikeLower := strings.ToLower(urlLike)
		urls = helpers.Filter(urls, func(url string) bool {
			return fuzzy.Match(urlLikeLower, strings.ToLower(url))
		})

		if len(urls) == 0 {
			return errors.New("no matches found in the given tap group")
		}
	}

	return br.OpenLinks(urls)
}
