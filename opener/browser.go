package opener

import (
	"os/exec"
	"runtime"
)

// browser implements the Opener interface and opens the links in the browser.
type browser struct {
}

// OpenLink opens a link in the browser
func (br *browser) OpenLink(request Request) error {
	link := request.Link

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
func (br *browser) OpenLinks(requests []Request) error {

	for _, r := range requests {
		err := br.OpenLink(r)
		if err != nil {
			return err
		}

	}
	return nil
}

func NewBrowser() Opener {
	return &browser{}
}
