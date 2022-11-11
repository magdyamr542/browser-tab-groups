package configManager

type ConfigManager interface {
	//GetConfig gets the config instance. this should be a map from a tap group to the list of urls
	GetConfig() (map[string][]string, error)

	// AddUrl adds the given url to the given tap group. the tap group is created if it does not exist
	AddUrl(url string, tapGroup string) error

	// GetUrls gets the urls for the given tap group
	GetUrls(tapGroup string) ([]string, error)

	// RemoveTapGroup removes all urls saved in the given tap group
	RemoveTapGroup(tapGroup string) error
}
