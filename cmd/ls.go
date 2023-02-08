package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/magdyamr542/browser-tab-groups/configManager"
)

// Listing all tap groups with their urls
type LsCmd struct{}

func (ls *LsCmd) Run() error {

	jsonCmg, err := configManager.NewJsonConfigManager()
	if err != nil {
		return err
	}
	return listTapGroups(os.Stdout, jsonCmg)
}

// listTapGroups lists all tap groups
func listTapGroups(outputW io.Writer, cm configManager.ConfigManager) error {
	cfg, err := cm.GetConfig()
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(cfg))
	for k := range cfg {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		groupName := key
		urls := cfg[key]
		entry := fmt.Sprintf("%v:\n", groupName)
		outputW.Write([]byte(entry))
		for i := range urls {
			outputW.Write([]byte(fmt.Sprintf("    %v\n", urls[i])))
		}
	}
	return nil
}
