package cmd

import (
	"io"
	"os"

	"github.com/magdyamr542/browser-tab-groups/configManager"
)

// Adding a new tap group
type AddCmd struct {
	TapGroup string `arg:"" name:"tap group" help:"the tap group to add the url to"`
	Url      string `arg:"" name:"url" help:"the url to add"`
}

func (add *AddCmd) Run() error {
	jsonCmg, err := configManager.NewJsonConfigManager()
	if err != nil {
		return err
	}
	return addUrlToTapGroup(os.Stdout, jsonCmg, add.TapGroup, add.Url)
}

// addUrlToTapGroup adds the given url to the given tap group
func addUrlToTapGroup(outputW io.Writer, cm configManager.ConfigManager, tapGroup, url string) error {
	err := cm.AddUrl(url, tapGroup)
	if err != nil {
		return err
	}
	outputW.Write([]byte("url added"))
	return nil
}
