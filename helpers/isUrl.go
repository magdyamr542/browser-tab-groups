package helpers

import "net/url"

// IsUrl checks if the given string s is a url
func IsUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil

}
