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
	errInputIsNotUrl          error = errors.New("the given input is not a url")
	errUrlIsAlreadyInTapGroup error = errors.New("the url is already in the tap group")
)

type Db map[string]any

// jsonConfigManager is an internal implementation for the config manager that saves the data as a json file
type jsonConfigManager struct {
	dirPath  string
	fileName string
	homeDir  string
}

func (cm *jsonConfigManager) GetConfig() (string, error) {

	var result strings.Builder

	if err := cm.walkDb(func(si scanInput) bool {
		result.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(Tap, len(si.parentGroups)), si.group))
		// Output the urls
		if si.isLeaf {
			for _, url := range si.value.([]string) {
				result.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(Tap, len(si.parentGroups)+1), url))
			}
		}

		return false
	}); err != nil {
		return "", err
	}

	return result.String(), nil
}
func (cm *jsonConfigManager) ExecForMatchingTapGroup(matcher func(tapGroupPath []string) bool, exec func(urls []string)) error {
	// e.g {"group.child.final" : {}}
	matchingPrefixes := make(map[string]struct{})
	leafs := make([]scanInput, 0)
	urls := make([]string, 0)

	if err := cm.walkDb(func(si scanInput) bool {
		path := append(si.parentGroups, si.group)
		if matcher(path) {
			matchingPrefixes[strings.Join(path, ".")] = struct{}{}
		}

		if si.isLeaf {
			leafs = append(leafs, si)
		}

		return false
	}); err != nil {
		return err
	}

	for _, leaf := range leafs {
		leafPath := append(leaf.parentGroups, leaf.group)
		leafPrefix := strings.Join(leafPath, ".")

		for matchingGroup := range matchingPrefixes {
			if strings.HasPrefix(leafPrefix, matchingGroup) {
				leafUrls := leaf.value.([]string)
				urls = append(urls, leafUrls...)
			}
		}
	}

	if len(urls) > 0 {
		exec(urls)
	}
	return nil
}

func (cm *jsonConfigManager) AddUrl(url string, tapGroups ...string) error {

	// Validate
	trimmedUrl := strings.TrimSpace(url)
	if !helpers.IsUrl(strings.TrimSpace(trimmedUrl)) {
		return errInputIsNotUrl
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
	if err := cm.walk(func(si scanInput) bool {
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

	file, _ := json.MarshalIndent(data, "", " ")

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

func NewJsonConfigManager() (ConfigManager, error) {
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

type scanInput struct {
	group        string
	value        any
	parentGroups []string
	isLeaf       bool // if true, value will be a slice of urls (slice of strings)
}

// the scanner should return false to indicate to stop the walking
type scanner func(scanInput) bool

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
			done = scanner(scanInput{group: currentGroup, value: urls, parentGroups: parentGroups, isLeaf: true})
			return done
		}

		// Nested
		nestedDb, isNested := value.(map[string]any)
		if isNested {
			// First visit the parent
			done = scanner(scanInput{group: currentGroup, value: value, parentGroups: parentGroups})
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

// func (cm *jsonConfigManager) printDb(db Db) {
// 	dbBytes, err := json.MarshalIndent(db, "", "  ")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%v\n", string(dbBytes))
// }

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
