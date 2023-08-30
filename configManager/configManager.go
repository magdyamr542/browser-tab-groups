package configManager

type ConfigManager interface {
	//GetConfig gets the config instance. this should be a map from a tap group to the list of urls
	GetConfig() (string, error)

	// AddUrl adds the given url to the given tap group. the tap group is created if it does not exist
	AddUrl(url string, tapGroup ...string) error

	// RemoveTapGroup removes all urls saved in the given tap group
	RemoveTapGroup(path ...string) error

	// ExecForMatchingTapGroup executes the functions (given the found urls) if the matcher returns true given
	// the current tapGroupPath
	ExecForMatchingTapGroup(matcher func(tapGroupPath []string) bool, exec func(urls []string)) error
}
