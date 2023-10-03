package opener

type Request struct {
	Link string
	// If not empty, will be logged before opening the link if the implementation is suitable to logging.
	Description string
}

// Opener opens links.
type Opener interface {
	//OpenLink opens a link in the browser
	OpenLink(Request) error

	//OpenLinks opens a link in the browser
	OpenLinks([]Request) error
}
