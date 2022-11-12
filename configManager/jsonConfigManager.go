package configManager

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/magdyamr542/browser-tab-groups/helpers"
)

// jsonConfigManager is an internal implementation for the config manager that saves the data as a json file
type jsonConfigManager struct {
	dirPath  string
	fileName string
	homeDir  string
}

var errGroupDoesNotExist error = errors.New("the group does not exist")
var errInputIsNotUrl error = errors.New("the given input is not a url")
var errUrlIsAlreadyInTapGroup error = errors.New("the url is already in the tap group")

func (cm *jsonConfigManager) GetConfig() (map[string][]string, error) {
	byteValue, _ := os.ReadFile(cm.storagePath())
	if len(byteValue) == 0 {
		byteValue = append(byteValue, []byte("{}")...)
	}
	var result map[string][]string
	err := json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return map[string][]string{}, err
	}
	return result, nil
}

func (cm *jsonConfigManager) GetTapGroups() ([]string, error) {
	db, err := cm.GetConfig()
	if err != nil {
		return []string{}, err
	}

	keys := make([]string, 0, len(db))
	for k := range db {
		keys = append(keys, k)
	}
	return keys, nil
}

func (cm *jsonConfigManager) GetUrls(tapGroup string) ([]string, error) {
	db, err := cm.GetConfig()
	if err != nil {
		return []string{}, err
	}

	trimmedTapGroup := strings.TrimSpace(tapGroup)
	_, ok := db[trimmedTapGroup]
	if !ok {
		return []string{}, errGroupDoesNotExist
	}

	return db[trimmedTapGroup], nil
}

func (cm *jsonConfigManager) AddUrl(url, tapGroup string) error {
	db, err := cm.GetConfig()
	if err != nil {
		return err
	}
	_, ok := db[tapGroup]
	if !ok {
		db[tapGroup] = []string{}
	}

	// Validate
	trimmedUrl := strings.TrimSpace(url)
	if !helpers.IsUrl(strings.TrimSpace(trimmedUrl)) {
		return errInputIsNotUrl
	}
	if helpers.Contains(db[tapGroup], strings.TrimSpace(trimmedUrl)) {
		return errUrlIsAlreadyInTapGroup
	}

	db[tapGroup] = append(db[tapGroup], strings.TrimSpace(trimmedUrl))
	return cm.refreshStorage(db)
}

func (cm *jsonConfigManager) RemoveTapGroup(tapGroup string) error {
	db, err := cm.GetConfig()
	if err != nil {
		return err
	}
	trimmedTapGroup := strings.TrimSpace(tapGroup)
	_, ok := db[trimmedTapGroup]
	if !ok {
		return errGroupDoesNotExist
	}
	delete(db, trimmedTapGroup)
	return cm.refreshStorage(db)
}

// overrides the saved file
func (cm *jsonConfigManager) refreshStorage(data map[string][]string) error {
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
