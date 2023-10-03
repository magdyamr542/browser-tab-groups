package configManager

type JsonConfigManager interface {
	ConfigManager

	//GetConfigJson returns the config as it's stored
	GetConfigJson() (string, error)

	//OverrideConfigJson overrides the stored config
	OverrideConfigJson(newConfig []byte) error
}

type TabGroup interface {
	// Returns all urls under the given TabGroup.
	Urls() ([]string, error)

	// Name of the current TabGroup.
	Name() string

	// Path to the current TabGroup
	// E.g [work, ticket1, github].
	Path() []string

	// Returns all children of the current TabGroup.
	Children() ([]TabGroup, error)

	// Formats the TabGroup as a string.
	String(prefix string) (string, error)
}

type ConfigManager interface {
	GetMatchingTabGroups(matcher func(tgPath []string) bool) ([]TabGroup, error)

	// AddUrl adds the given url to the given tab group. the tab group is created if it does not exist
	AddUrl(url string, tabGroup ...string) error

	// RemoveTabGroup removes all urls saved in the given tab group
	RemoveTabGroup(path ...string) error
}
