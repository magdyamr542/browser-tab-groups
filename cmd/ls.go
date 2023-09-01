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
	TapGroups []string `arg:"" optional:"" name:"path to tap group" help:"the path to the tap group to remove"`
}

func (ls *LsCmd) Run() error {

	jsonCmg, err := configManager.NewJsonConfigManager()
	if err != nil {
		return err
	}
	return listTapGroups(os.Stdout, jsonCmg, ls.TapGroups)
}

// listTapGroups lists all tap groups
func listTapGroups(outputW io.Writer, cm configManager.ConfigManager, tapGroups []string) error {
	if len(tapGroups) == 0 {
		cfg, err := cm.GetConfig()
		if err != nil {
			return err
		}
		outputW.Write([]byte(cfg))
		return nil
	}

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

	for _, url := range urls {
		outputW.Write([]byte(url + "\n"))
	}

	return nil
}
