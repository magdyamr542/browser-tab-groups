package cmd

import (
	"io"
	"os"

	"github.com/magdyamr542/browser-tab-groups/configManager"
)

// Removing a tap group
type RmCmd struct {
	TapGroup string `arg:"" name:"tap group" help:"the tap group to remove"`
}

func (rm *RmCmd) Run() error {

	jsonCmg, err := configManager.NewJsonConfigManager()
	if err != nil {
		return err
	}
	return removeTapGroup(os.Stdout, jsonCmg, rm.TapGroup)
}

// removeTapGroup removes a saved tap group
func removeTapGroup(outputW io.Writer, cm configManager.ConfigManager, tapGroup string) error {
	err := cm.RemoveTapGroup(tapGroup)
	if err != nil {
		return err
	}
	outputW.Write([]byte("removed"))
	return nil
}
