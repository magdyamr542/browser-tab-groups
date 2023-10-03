package opener

import (
	"fmt"
	"io"
	"net/http"
)

var (
	newline = []byte("\n")
)

// shell implements the Opener interface and opens the links in the shell.
type shell struct {
	output io.Writer
}

// OpenLink opens a link in the browser
func (s *shell) OpenLink(request Request) error {
	fmt.Printf(request.Description)

	link := request.Link

	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body = append(body, newline...)

	if _, err := s.output.Write(body); err != nil {
		return err
	}

	return nil
}

// OpenLinks opens all links in the browser
func (s *shell) OpenLinks(requests []Request) error {

	for _, r := range requests {
		err := s.OpenLink(r)
		if err != nil {
			return err
		}

	}
	return nil
}

func NewShell(output io.Writer) Opener {
	return &shell{output: output}
}
