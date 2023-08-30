package cmd

import (
	"io"
	"os"

	"github.com/magdyamr542/browser-tab-groups/configManager"
)

// Removing a tap group
type RmCmd struct {
	TapGroups []string `arg:"" name:"path to tap group" help:"the path to the tap group to remove"`
}

func (rm *RmCmd) Run() error {

	jsonCmg, err := configManager.NewJsonConfigManager()
	if err != nil {
		return err
	}
	return removeTapGroup(os.Stdout, jsonCmg, rm.TapGroups...)
}

// removeTapGroup removes a saved tap group
func removeTapGroup(outputW io.Writer, cm configManager.ConfigManager, path ...string) error {
	err := cm.RemoveTapGroup(path...)
	if err != nil {
		return err
	}
	outputW.Write([]byte("removed"))
	return nil
}
