package configManager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/magdyamr542/browser-tab-groups/helpers"
)

const Tap = "  "

var (
	errUrlIsAlreadyInTapGroup error = errors.New("the url is already in the tap group")
)

type Db map[string]any

// jsonConfigManager is an internal implementation for the config manager that saves the data as a json file
type jsonConfigManager struct {
	dirPath  string
	fileName string
	homeDir  string
}

func (cm *jsonConfigManager) GetMatchingTapGroups(matcher func(tapGroupPath []string) bool) ([]TapGroup, error) {

	tgs := make([]TapGroup, 0)

	onMatch := func(tg TapGroup) {
		tgs = append(tgs, tg)
	}

	err := cm.ExecForMatchingTapGroup(matcher, onMatch)

	if err != nil {
		return nil, err
	}

	if len(tgs) == 0 {
		return nil, fmt.Errorf("no matching tap groups found")
	}

	return tgs, nil

}

func (cm *jsonConfigManager) OverrideConfigJson(newConfig []byte) error {
	var newDb Db
	if err := json.Unmarshal(newConfig, &newDb); err != nil {
		return err
	}
	return cm.refreshStorage(newDb)
}

func (cm *jsonConfigManager) GetConfigJson() (string, error) {

	db, err := cm.getDB()
	if err != nil {
		return "", err
	}

	byteValue, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return "", err
	}

	return string(byteValue), nil
}

func (cm *jsonConfigManager) ExecForMatchingTapGroup(matcher func(tapGroupPath []string) bool, exec func(TapGroup)) error {

	db, err := cm.getDB()
	if err != nil {
		return err
	}

	var retErr error

	if err := cm.walk(func(si entry) bool {
		if matcher(si.Path()) {

			tg, err := NewTg(&db, si.Path())
			if err != nil {
				retErr = err
				return true
			}

			exec(tg)
		}

		return false
	}, db); err != nil {
		return err
	}

	return retErr
}

func (cm *jsonConfigManager) AddUrl(url string, tapGroups ...string) error {

	// Validate
	trimmedUrl := strings.TrimSpace(url)
	if !helpers.IsUrl(strings.TrimSpace(trimmedUrl)) {
		return fmt.Errorf("%q is not a url", url)
	}

	db, err := cm.getDB()
	if err != nil {
		return err
	}

	// Create the nested tap groups if necessary and add the url to the leaf
	currentTapGroup := 0
	currentDb := db
	for currentTapGroup < len(tapGroups) {
		tapGroup := tapGroups[currentTapGroup]
		_, ok := currentDb[tapGroup]
		if !ok {
			// Last tap group. This maps to the list of urls
			if currentTapGroup+1 >= len(tapGroups) {
				currentDb[tapGroup] = []string{url}
			} else {
				// Go deeper
				currentDb[tapGroup] = make(Db)
				currentDb = currentDb[tapGroup].(Db)
			}
		} else {
			// Key exists
			urlsAny, isLeaf := currentDb[tapGroup].([]any)
			if isLeaf {
				// User trying to create a new tap group under an existing leaf. Error
				if currentTapGroup != len(tapGroups)-1 {
					return fmt.Errorf("can't create %q as a tap group inside %[2]q (%[2]q already contains urls)",
						tapGroups[len(tapGroups)-1], tapGroup)
				}

				// Add the url to the existing urls
				currentUrls := getUrls(urlsAny)
				if helpers.Contains(currentUrls, trimmedUrl) {
					return errUrlIsAlreadyInTapGroup
				}
				currentUrls = append(currentUrls, url)
				currentDb[tapGroup] = currentUrls

			} else {
				// Go one level deeper
				currentDb = currentDb[tapGroup].(map[string]any)
			}
		}
		currentTapGroup += 1
	}

	// cm.printDb(db)

	return cm.refreshStorage(db)
}

func (cm *jsonConfigManager) RemoveTapGroup(path ...string) error {

	db, err := cm.getDB()
	if err != nil {
		return err
	}

	found := false
	if err := cm.walk(func(si entry) bool {
		if equal(append(si.parentGroups, si.group), path) {
			found = true
			// Delete it
			currentDb := db
			for _, parentGroup := range si.parentGroups {
				currentDb = currentDb[parentGroup].(map[string]any)
			}
			delete(currentDb, si.group)
			return true
		}
		return false
	}, db); err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("no match found for %s", path)
	}

	// cm.printDb(db)

	return cm.refreshStorage(db)
}

// overrides the saved file
func (cm *jsonConfigManager) refreshStorage(data Db) error {
	storagePath := cm.storagePath()
	flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	f, err := os.OpenFile(storagePath, flags, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(storagePath, file, 0600)

}

// create dir and file for storage
func (cm *jsonConfigManager) initStorage() error {

	err := os.MkdirAll(path.Join(cm.homeDir, cm.dirPath), os.ModePerm)
	if err != nil {
		return err
	}
	_, err = os.OpenFile(cm.storagePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (cm *jsonConfigManager) storagePath() string {
	return filepath.Join(cm.homeDir, cm.dirPath, cm.fileName)
}

func NewJsonConfigManager() (JsonConfigManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cm := jsonConfigManager{dirPath: ".browser-tabs", fileName: "tabs.json", homeDir: homeDir}
	err = cm.initStorage()
	if err != nil {
		return nil, err
	}
	return &cm, nil
}

type entry struct {
	group        string
	value        any
	parentGroups []string // if nil or has 0 length, this is the root
	isLeaf       bool     // if true, value will be a slice of urls (slice of strings)
}

func (s *entry) PathStr() string {
	return strings.Join(s.Path(), ".")
}

func (s *entry) Path() []string {
	return append(s.parentGroups, s.group)
}

func (s *entry) Root() bool {
	return len(s.parentGroups) == 0
}

func (s *entry) EmptySpacePrefix() string {
	return strings.Repeat(Tap, len(s.parentGroups))
}

// the scanner should return false to indicate to stop the walking
type scanner func(entry) bool

// walk walks the JSON giving the scanner access to the different groups
func (cm *jsonConfigManager) walk(scanner scanner, db Db) error {
	done := false

	var walkRecursive func(string, []string, any) bool
	walkRecursive = func(currentGroup string, parentGroups []string, value any) bool {
		if done {
			return true
		}

		// Leaf
		urlsAny, isLeaf := value.([]any)
		if isLeaf {
			urls := getUrls(urlsAny)
			done = scanner(entry{group: currentGroup, value: urls, parentGroups: parentGroups, isLeaf: true})
			return done
		}

		// Nested
		nestedDb, isNested := value.(map[string]any)
		if isNested {
			// First visit the parent
			done = scanner(entry{group: currentGroup, value: value, parentGroups: parentGroups})
			if done {
				return true
			}
			// Then visit the children
			childGroups := make([]string, 0)
			for childGroup := range nestedDb {
				childGroups = append(childGroups, childGroup)
			}
			sort.Strings(childGroups)
			for _, childGroup := range childGroups {
				newParentGroups := make([]string, len(parentGroups))
				copy(newParentGroups, parentGroups)
				if walkRecursive(childGroup, append(newParentGroups, currentGroup), nestedDb[childGroup]) {
					return true
				}
			}

			return false
		}

		return true
	}

	// Trigger the walking
	groups := make([]string, 0)
	for key := range db {
		groups = append(groups, key)
	}
	sort.Strings(groups)
	for _, group := range groups {
		if walkRecursive(group, []string{}, db[group]) {
			return nil
		}
	}

	return nil
}

func (cm *jsonConfigManager) walkDb(scanner scanner) error {
	db, err := cm.getDB()
	if err != nil {
		return err
	}

	return cm.walk(scanner, db)
}

func (cm *jsonConfigManager) getDB() (Db, error) {
	byteValue, err := os.ReadFile(cm.storagePath())
	if err != nil {
		return nil, err
	}

	if len(byteValue) == 0 {
		byteValue = append(byteValue, []byte("{}")...)
	}

	var result Db
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getUrls(urlsAny []any) []string {
	urls := make([]string, 0, len(urlsAny))
	for _, u := range urlsAny {
		urls = append(urls, u.(string))
	}
	return urls
}

func equal(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

type tg struct {
	// The root db.
	db *Db

	// path to the current TapGroup.
	path []string

	urls []string
}

func (t *tg) Leaf() bool {
	return len(t.urls) != 0
}

func NewTg(db *Db, path []string) (*tg, error) {
	currDb := *db
	tg := tg{db: db, path: path}

	for i, key := range path {
		nestedDb, ok := currDb[key]
		if !ok {
			return nil, fmt.Errorf("no such tap group: %s", strings.Join(path, "."))
		}

		urlsAny, isLeaf := nestedDb.([]any)
		if isLeaf && i != len(path)-1 {
			return nil, fmt.Errorf("no such tap group: %s", strings.Join(path, "."))
		}

		if isLeaf {
			tg.urls = getUrls(urlsAny)
		} else {

			currDb = nestedDb.(map[string]any)
		}

	}
	return &tg, nil
}

// tg implements the TapGroup interface.
// Returns all urls under the given TapGroup.
func (tg *tg) Urls() ([]string, error) {
	if tg.Leaf() {
		return tg.urls, nil
	}

	children, err := tg.Children()
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for i := range children {
		childUrls, err := children[i].Urls()
		if err != nil {
			return nil, err
		}
		urls = append(urls, childUrls...)
	}

	return urls, nil
}

// Name of the current TapGroup.
func (tg *tg) Name() string {
	if len(tg.path) == 0 {
		return ""
	}
	return tg.path[len(tg.path)-1]
}

// Path to the current TapGroup
// E.g [work, ticket1, github].
func (tg *tg) Path() []string {
	return tg.path
}

// Returns all children of the current TapGroup.
func (t *tg) Children() ([]TapGroup, error) {
	if t.Leaf() {
		return nil, nil
	}

	currDb := *t.db
	for _, key := range t.path {
		currDb = currDb[key].(map[string]any)
	}

	result := make([]TapGroup, 0)
	for key := range currDb {
		tg, err := NewTg(t.db, append(t.path, key))
		if err != nil {
			return nil, err
		}
		result = append(result, tg)
	}

	return result, nil
}

// Formats the TapGroup as a string.
func (tg *tg) String(prefix string) (string, error) {
	var writer strings.Builder

	writer.WriteString(fmt.Sprintf("%s%s\n", prefix, asBold(tg.Name())))

	if tg.Leaf() {
		for i, url := range tg.urls {
			writer.WriteString(fmt.Sprintf("%s%s%s", prefix, Tap, url))
			if i != len(tg.urls)-1 {
				writer.WriteString("\n")
			}
		}
		return writer.String(), nil
	}

	children, err := tg.Children()
	if err != nil {
		return "", err
	}
	for i := range children {
		childStr, err := children[i].String(prefix + Tap)
		if err != nil {
			return "", err
		}
		writer.WriteString(childStr)

		if i != len(children)-1 {
			writer.WriteString("\n")
		}
	}
	return writer.String(), nil
}

func asBold(str string) string {
	return "\x1b[1m" + str + "\x1b[0m"
}
