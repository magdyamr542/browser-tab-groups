package helpers

// Contains checks if a slice s contains a string
func Contains(s []string, str string) bool {
	for _, a := range s {
		if a == str {
			return true
		}
	}
	return false
}
