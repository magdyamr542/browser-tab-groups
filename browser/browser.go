package browser

import (
	"os/exec"
	"runtime"
)

type Browser interface {
	//OpenLink opens a link in the browser
	OpenLink(link string) error

	//OpenLinks opens a link in the browser
	OpenLinks(links []string) error
}

// browser is the internal implementation for the Browser interface
type browser struct {
}

// OpenLink opens a link in the browser
func (br *browser) OpenLink(link string) error {

	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], link)...)
	return cmd.Start()
}

// OpenLinks opens all links in the browser
func (br *browser) OpenLinks(links []string) error {

	for _, link := range links {
		err := br.OpenLink(link)
		if err != nil {
			return err
		}

	}
	return nil
}

func NewBrowser() Browser {
	return &browser{}
}
