package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/magdyamr542/browser-tab-groups/configManager"
	"github.com/urfave/cli/v2"
)

var RmCmd cli.Command = cli.Command{
	Name:        "remove",
	Usage:       "Remove a saved tab group",
	Description: "Remove a saved tab group using its path",
	Aliases:     []string{"rm"},
	Action: func(cCtx *cli.Context) error {
		jsonCmg, err := configManager.NewJsonConfigManager()
		if err != nil {
			return err
		}

		tabGroups := cCtx.Args().Slice()
		if len(tabGroups) == 0 {
			return fmt.Errorf("provide a path to the tab group you want to delete (as space separated string)")
		}
		return removeTabGroup(os.Stdout, jsonCmg, tabGroups...)
	},
}

// removeTabGroup removes a saved tab group
func removeTabGroup(outputW io.Writer, cm configManager.ConfigManager, path ...string) error {
	err := cm.RemoveTabGroup(path...)
	if err != nil {
		return err
	}
	outputW.Write([]byte("removed"))
	return nil
}
