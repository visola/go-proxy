package upstream

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/visola/go-proxy/pkg/configuration"
	"gopkg.in/yaml.v2"
)

const currentCustomDirectoriesFile = "custom_directories.yml"

var (
	customDirectories      = make([]string, 0)
	customDirectoriesMutex sync.Mutex
)

// AddCustomDirectory adds a directory to the current list of directories
func AddCustomDirectory(toAdd string) {
	for _, d := range customDirectories {
		if toAdd == d {
			return
		}
	}

	saveCustomDirectories(append(customDirectories, toAdd))
}

// CustomDirectories return all the custom directories configured
func CustomDirectories() []string {
	return customDirectories
}

// RemoveCustomDirectory removes a specific directory from the list
func RemoveCustomDirectory(toRemove string) {
	newState := make([]string, len(customDirectories)-1)
	j := 0
	for _, c := range customDirectories {
		if c != toRemove {
			newState[j] = c
			j++
		}
	}
	saveCustomDirectories(newState)
}

func loadCustomDirectories() error {
	configDir, configDirErr := configuration.GetConfigurationDirectory()
	if configDirErr != nil {
		return configDirErr
	}

	persistedStateFile := filepath.Join(configDir, currentCustomDirectoriesFile)

	yamlContent, readErr := ioutil.ReadFile(persistedStateFile)
	if readErr != nil {
		return readErr
	}

	var loadedCustomDirectories []string
	if unmarshalError := yaml.Unmarshal(yamlContent, &loadedCustomDirectories); unmarshalError != nil {
		return unmarshalError
	}

	customDirectoriesMutex.Lock()
	defer customDirectoriesMutex.Unlock()

	customDirectories = loadedCustomDirectories
	return nil
}

func saveCustomDirectories(toSave []string) error {
	customDirectoriesMutex.Lock()
	defer customDirectoriesMutex.Unlock()

	configDir, err := configuration.GetConfigurationDirectory()
	if err != nil {
		return err
	}

	persistedStateFile := filepath.Join(configDir, currentCustomDirectoriesFile)

	customDirectories = toSave
	data, marshalErr := yaml.Marshal(customDirectories)
	if marshalErr != nil {
		return marshalErr
	}

	return ioutil.WriteFile(persistedStateFile, data, 0644)
}
