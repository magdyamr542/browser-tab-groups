package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/magdyamr542/browser-tab-groups/configManager"
)

// Adding a new tap group
type AddCmd struct {
	GroupsThenUrl []string `arg:"" name:"tap groups then the url" help:"the tap groups to add the url to. groups are hierarchical"`
}

func (add *AddCmd) Run() error {

	example := `1. browser-tab-group add work issue1 https://www.google.com (url is nested under 2 groups)
2. browser-tab-group add work https://www.google.com (url is nested under 1 group)
`
	if len(add.GroupsThenUrl) == 1 {
		return fmt.Errorf("wrong usage. Needs at least one group and a url (2 inputs)\n%s", example)
	}
	url := add.GroupsThenUrl[len(add.GroupsThenUrl)-1]
	groups := add.GroupsThenUrl[:len(add.GroupsThenUrl)-1]
	jsonCmg, err := configManager.NewJsonConfigManager()
	if err != nil {
		return err
	}
	return addUrlToTapGroup(os.Stdout, jsonCmg, url, groups...)
}

func addUrlToTapGroup(outputW io.Writer, cm configManager.ConfigManager, url string, tapGroups ...string) error {
	err := cm.AddUrl(url, tapGroups...)
	if err != nil {
		return err
	}
	outputW.Write([]byte("url added"))
	return nil
}
