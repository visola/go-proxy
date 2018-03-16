package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

const stateFileName = ".current_state"

var configurations []DynamicMapping

// GetConfigurations get all configurations or load them if not loaded so far
func GetConfigurations() ([]DynamicMapping, error) {
	if configurations != nil {
		return configurations, nil
	}

	var err error
	configurations, err = getCurrentState()
	return configurations, err
}

// SetStatus change the status of a configuration and save the new state to
// disk so that it will get loaded on next start
func SetStatus(mappingID string, status bool) error {
	_, err := GetConfigurations()
	if err != nil {
		return err
	}

	// TODO - Make sure this change is atomic
	for index, config := range configurations {
		if config.MappingID == mappingID {
			configurations[index].Active = status
			break
		}
	}
	return storeCurrentState()
}

func getConfigDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.go-proxy", nil
}

func getCurrentState() ([]DynamicMapping, error) {
	staticConfigs, staticConfigErr := loadAllConfigurations()
	if staticConfigErr != nil {
		return nil, staticConfigErr
	}

	storedStates, storedStateErr := getStoredState()
	if storedStateErr != nil {
		return nil, storedStateErr
	}

	result := make([]DynamicMapping, len(staticConfigs))
	for index, staticConfig := range staticConfigs {
		if storedState, exists := storedStates[staticConfig.MappingID]; exists {
			result[index] = storedState
		} else {
			result[index] = DynamicMapping{
				Active:    true,
				From:      staticConfig.From,
				MappingID: staticConfig.MappingID,
				Origin:    staticConfig.Origin,
				Proxy:     staticConfig.Proxy,
				Regexp:    staticConfig.Regexp,
				To:        staticConfig.To,
			}
		}
	}

	return result, nil
}

func getStoredState() (map[string]DynamicMapping, error) {
	configDir, configDirErr := getConfigDirectory()
	if configDirErr != nil {
		return nil, configDirErr
	}

	currentStateFile, curretnStateFileErr := os.Open(path.Join(configDir, stateFileName))
	if os.IsNotExist(curretnStateFileErr) {
		return make(map[string]DynamicMapping, 0), nil
	}

	if curretnStateFileErr != nil {
		return nil, curretnStateFileErr
	}

	currentStateData, readErr := ioutil.ReadAll(currentStateFile)
	if readErr != nil {
		return nil, readErr
	}

	var result map[string]DynamicMapping
	unmarshalErr := json.Unmarshal(currentStateData, &result)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return result, nil
}

func loadAllConfigurations() ([]Mapping, error) {
	configDir, configDirErr := getConfigDirectory()
	if configDirErr != nil {
		panic(configDirErr)
	}

	files, filesErr := ioutil.ReadDir(configDir)
	if filesErr != nil {
		panic(filesErr)
	}

	configurations := make([]Mapping, 0)
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".yml" && filepath.Ext(file.Name()) != ".yaml" {
			continue
		}

		loadedConfigurations, loadingErr := loadConfigurations(configDir, file)
		if loadingErr != nil {
			return nil, loadingErr
		}
		configurations = append(configurations, loadedConfigurations...)
	}

	sortConfigurations(configurations)
	return configurations, nil
}

func loadConfigurations(configDir string, file os.FileInfo) ([]Mapping, error) {
	configurations := make([]Mapping, 0)

	config, configErr := readConfiguration(path.Join(configDir, file.Name()))
	if configErr != nil {
		return nil, configErr
	}

	for _, staticConfig := range config.Static {
		configurations = append(configurations, fromYAMLMapping(staticConfig, file.Name(), false))
	}

	for _, staticConfig := range config.Proxy {
		configurations = append(configurations, fromYAMLMapping(staticConfig, file.Name(), true))
	}

	return configurations, nil
}

func readConfiguration(file string) (loadedConfig proxyConfig, err error) {
	var yamlContent []byte
	yamlContent, err = ioutil.ReadFile(file)

	if err != nil {
		return loadedConfig, err
	}

	err = yaml.Unmarshal(yamlContent, &loadedConfig)
	return loadedConfig, err
}

func sortConfigurations(configurations []Mapping) {
	sort.Slice(configurations, func(i, j int) bool {
		pathI := strings.ToLower(configurations[i].From)
		pathJ := strings.ToLower(configurations[j].From)

		if pathI == "" {
			pathI = strings.ToLower(configurations[i].Regexp)
		}

		if pathJ == "" {
			pathJ = strings.ToLower(configurations[j].Regexp)
		}

		if len(pathI) == len(pathJ) {
			return strings.Compare(pathI, pathJ) < 0
		}

		return len(pathI) > len(pathJ)
	})
}

func storeCurrentState() error {
	toStore := make(map[string]DynamicMapping)
	for _, config := range configurations {
		toStore[config.MappingID] = config
	}

	// Force reload
	configurations = nil

	data, dataErr := json.Marshal(toStore)
	if dataErr != nil {
		return dataErr
	}

	configDir, configDirErr := getConfigDirectory()
	if configDirErr != nil {
		return configDirErr
	}

	return ioutil.WriteFile(path.Join(configDir, stateFileName), data, 0644)
}
