package cmd

import (
	"github.com/magdyamr542/browser-tab-groups/configManager"
	"github.com/magdyamr542/browser-tab-groups/editor"
)

// Editing the storage manually
type EditCmd struct {
}

func (ls *EditCmd) Run() error {
	jsonCmg, err := configManager.NewJsonConfigManager()
	if err != nil {
		return err
	}
	return editTapGroups(jsonCmg)
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
