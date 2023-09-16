package configManager

type JsonConfigManager interface {
	ConfigManager

	//GetConfigJson returns the config as it's stored
	GetConfigJson() (string, error)

	//OverrideConfigJson overrides the stored config
	OverrideConfigJson(newConfig []byte) error
}

type TapGroup interface {
	// Returns all urls under the given TapGroup.
	Urls() ([]string, error)

	// Name of the current TapGroup.
	Name() string

	// Path to the current TapGroup
	// E.g [work, ticket1, github].
	Path() []string

	// Returns all children of the current TapGroup.
	Children() ([]TapGroup, error)

	// Formats the TapGroup as a string.
	String(prefix string) (string, error)
}

type ConfigManager interface {
	GetMatchingTapGroups(matcher func(tapGroupPath []string) bool) ([]TapGroup, error)

	// AddUrl adds the given url to the given tap group. the tap group is created if it does not exist
	AddUrl(url string, tapGroup ...string) error

	// RemoveTapGroup removes all urls saved in the given tap group
	RemoveTapGroup(path ...string) error
}
