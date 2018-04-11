package mapping

import (
	"encoding/json"
	"fmt"
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

var mappings []DynamicMapping

// GetMappings get all mappings or load them if not loaded so far
func GetMappings() ([]DynamicMapping, error) {
	if mappings != nil {
		return mappings, nil
	}

	var err error
	mappings, err = getCurrentState()
	return mappings, err
}

// SetStatus change the status of a mapping and save the new state to
// disk so that it will get loaded on next start
func SetStatus(mappingID string, status bool) error {
	_, err := GetMappings()
	if err != nil {
		return err
	}

	// TODO - Make sure this change is atomic
	for index, mapping := range mappings {
		if mapping.MappingID == mappingID {
			mappings[index].Active = status
			break
		}
	}
	return storeCurrentState()
}

func getMappingDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.go-proxy", nil
}

func getCurrentState() ([]DynamicMapping, error) {
	staticMappings, staticMappingErr := loadAllMappings()
	if staticMappingErr != nil {
		return nil, staticMappingErr
	}

	storedStates, storedStateErr := getStoredState()
	if storedStateErr != nil {
		return nil, storedStateErr
	}

	result := make([]DynamicMapping, len(staticMappings))
	for index, staticMapping := range staticMappings {
		if storedState, exists := storedStates[staticMapping.MappingID]; exists {
			result[index] = storedState
		} else {
			result[index] = DynamicMapping{
				Active:    true,
				From:      staticMapping.From,
				MappingID: staticMapping.MappingID,
				Origin:    staticMapping.Origin,
				Proxy:     staticMapping.Proxy,
				Regexp:    staticMapping.Regexp,
				To:        staticMapping.To,
			}
		}

		validationErr := result[index].Validate()
		if validationErr != nil {
			return nil, fmt.Errorf("Error in mapping from file: %s\n%s", result[index].Origin, validationErr)
		}
	}

	return result, nil
}

func getStoredState() (map[string]DynamicMapping, error) {
	mappingDir, mappingDirErr := getMappingDirectory()
	if mappingDirErr != nil {
		return nil, mappingDirErr
	}

	currentStateFile, currentStateFileErr := os.Open(path.Join(mappingDir, stateFileName))
	if os.IsNotExist(currentStateFileErr) {
		return make(map[string]DynamicMapping, 0), nil
	}

	if currentStateFileErr != nil {
		return nil, currentStateFileErr
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

func loadAllMappings() ([]Mapping, error) {
	mappingDir, mappingDirErr := getMappingDirectory()
	if mappingDirErr != nil {
		return nil, mappingDirErr
	}

	files, filesErr := ioutil.ReadDir(mappingDir)
	if filesErr != nil {
		return nil, filesErr
	}

	mappings := make([]Mapping, 0)
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".yml" && filepath.Ext(file.Name()) != ".yaml" {
			continue
		}

		loadedMappings, loadingErr := loadMappings(mappingDir, file)
		if loadingErr != nil {
			return nil, loadingErr
		}
		mappings = append(mappings, loadedMappings...)
	}

	sortMappings(mappings)
	return mappings, nil
}

func loadMappings(mappingDir string, file os.FileInfo) ([]Mapping, error) {
	mappings := make([]Mapping, 0)

	mapping, mappingErr := readMapping(path.Join(mappingDir, file.Name()))
	if mappingErr != nil {
		return nil, fmt.Errorf("Error while reading configuration file: %s\n%s", file.Name(), mappingErr)
	}

	for _, staticMapping := range mapping.Static {
		mappings = append(mappings, fromYAMLMapping(staticMapping, file.Name(), false))
	}

	for _, staticMapping := range mapping.Proxy {
		mappings = append(mappings, fromYAMLMapping(staticMapping, file.Name(), true))
	}

	return mappings, nil
}

func readMapping(file string) (loadedMapping yamlMapping, err error) {
	var yamlContent []byte
	yamlContent, err = ioutil.ReadFile(file)

	if err != nil {
		return loadedMapping, err
	}

	err = yaml.Unmarshal(yamlContent, &loadedMapping)
	return loadedMapping, err
}

func sortMappings(mappings []Mapping) {
	sort.Slice(mappings, func(i, j int) bool {
		pathI := strings.ToLower(mappings[i].From)
		pathJ := strings.ToLower(mappings[j].From)

		if pathI == "" {
			pathI = strings.ToLower(mappings[i].Regexp)
		}

		if pathJ == "" {
			pathJ = strings.ToLower(mappings[j].Regexp)
		}

		if len(pathI) == len(pathJ) {
			return strings.Compare(pathI, pathJ) < 0
		}

		return len(pathI) > len(pathJ)
	})
}

func storeCurrentState() error {
	toStore := make(map[string]DynamicMapping)
	for _, mapping := range mappings {
		toStore[mapping.MappingID] = mapping
	}

	// Force reload
	mappings = nil

	data, dataErr := json.Marshal(toStore)
	if dataErr != nil {
		return dataErr
	}

	mappingDir, mappingDirErr := getMappingDirectory()
	if mappingDirErr != nil {
		return mappingDirErr
	}

	return ioutil.WriteFile(path.Join(mappingDir, stateFileName), data, 0644)
}
