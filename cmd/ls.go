package cmd

import (
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
	outputW.Write([]byte(cfg))

	return nil
}
