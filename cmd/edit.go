package cmd

import (
	"github.com/magdyamr542/browser-tab-groups/configManager"
	"github.com/magdyamr542/browser-tab-groups/editor"
	"github.com/urfave/cli/v2"
)

// Editing the storage manually
var EditCmd cli.Command = cli.Command{
	Name:        "edit",
	Usage:       "Edit tab groups manually using your editor",
	Description: "Edit the tap groups manually (JSON editing in you editor)",
	Action: func(cCtx *cli.Context) error {

		jsonCmg, err := configManager.NewJsonConfigManager()
		if err != nil {
			return err
		}
		return editTapGroups(jsonCmg)
	},
}

func editTapGroups(cm configManager.JsonConfigManager) error {

	cfgJson, err := cm.GetConfigJson()
	if err != nil {
		return err
	}

	newContent, err := editor.New().Edit([]byte(cfgJson))
	if err != nil {
		return err
	}

	return cm.OverrideConfigJson(newContent)
}
