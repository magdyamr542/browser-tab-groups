package cmd

import (
	"fmt"
	"io"
	"os"

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

	i := 0
	for groupName, urls := range cfg {
		entry := fmt.Sprintf("%v:\n", groupName)
		outputW.Write([]byte(entry))
		for i := range urls {
			outputW.Write([]byte(fmt.Sprintf(" %v\n", urls[i])))
		}
		i += 1
		if i < len(cfg) {
			outputW.Write([]byte("\n"))
		}
	}
	return nil
}
