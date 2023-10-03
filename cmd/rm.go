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
	Usage:       "Remove a saved tap group",
	Description: "Remove a saved tap group using its path",
	Aliases:     []string{"rm"},
	Action: func(cCtx *cli.Context) error {
		jsonCmg, err := configManager.NewJsonConfigManager()
		if err != nil {
			return err
		}

		tapGroups := cCtx.Args().Slice()
		if len(tapGroups) == 0 {
			return fmt.Errorf("provide a path to the tap group you want to delete (as space separated string)")
		}
		return removeTapGroup(os.Stdout, jsonCmg, tapGroups...)
	},
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
