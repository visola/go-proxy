package mapping

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path"
	"path/filepath"
)

const stateFileName = ".current_state"

var mappings []Mapping

// GetMappings get all mappings or load them if not loaded so far
func GetMappings() ([]Mapping, error) {
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

func getCurrentState() ([]Mapping, error) {
	mappingsFromFiles, mappingsFromFilesErr := loadAllMappings()
	if mappingsFromFilesErr != nil {
		return nil, mappingsFromFilesErr
	}

	mappingsFromStore, mappingsFromStoreErr := getStoredState()
	if mappingsFromStoreErr != nil {
		return nil, mappingsFromStoreErr
	}

	result := make([]Mapping, len(mappingsFromFiles))
	for index, mappingFromFiles := range mappingsFromFiles {
		if mappingFromStore, exists := mappingsFromStore[mappingFromFiles.MappingID]; exists {
			result[index] = mappingFromStore
		} else {
			result[index] = mappingFromFiles
		}

		validationErr := result[index].Validate()
		if validationErr != nil {
			return nil, fmt.Errorf("Error in mapping from file: %s\n%s", result[index].Origin, validationErr)
		}
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

		loadedMappings, loadingErr := loadMappingsFromFiles(mappingDir, file)
		if loadingErr != nil {
			return nil, loadingErr
		}
		mappings = append(mappings, loadedMappings...)
	}

	sortMappings(mappings)
	return mappings, nil
}

func storeCurrentState() error {
	toStore := make(map[string]Mapping)
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
