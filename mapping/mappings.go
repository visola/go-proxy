package mapping

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/Everbridge/go-proxy/configuration"
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

// Set sets a mapping from its previous value to a new value
func Set(mappingID string, mapping Mapping) error {
	_, err := GetMappings()
	if err != nil {
		return err
	}

	mapping.MappingID = generateID(mapping)

	// TODO - Make sure this change is atomic
	for index, toBeReplaced := range mappings {
		if toBeReplaced.MappingID == mappingID {
			mappings[index] = mapping
			break
		}
	}
	return storeCurrentState()
}

// SetAll sets all the mappings to a current new state
func SetAll(newMappings []Mapping) error {
	_, err := GetMappings()
	if err != nil {
		return err
	}

	// Replace all
	mappings = newMappings
	return storeCurrentState()
}

func findAllMappingsFiles() ([]string, error) {
	result := make([]string, 0)
	filesFromBaseDirs, err := findFilesFromBaseDirectories()
	if err != nil {
		return nil, err
	}

	result = append(result, filesFromBaseDirs...)

	filesFromConfigDir, err := findFilesFromConfigurationDirector()
	if err != nil {
		return nil, err
	}

	result = append(result, filesFromConfigDir...)
	return result, nil
}

func findFilesFromBaseDirectories() ([]string, error) {
	result := make([]string, 0)

	config, loadConfigErr := configuration.LoadConfiguration()
	if loadConfigErr != nil {
		return nil, loadConfigErr
	}

	walkDirectories(config.BaseDirectories, func(pathToFile string, isDir bool) {
		if !isDir && isGoProxyMappingFile(pathToFile) {
			result = append(result, pathToFile)
		}
	})

	return result, nil
}

func findFilesFromConfigurationDirector() ([]string, error) {
	result := make([]string, 0)

	mappingDir, mappingDirErr := configuration.GetConfigurationDirectory()
	if mappingDirErr != nil {
		return nil, mappingDirErr
	}

	files, filesErr := ioutil.ReadDir(mappingDir)
	if filesErr != nil {
		return nil, filesErr
	}

	for _, file := range files {
		pathToFile := path.Join(mappingDir, file.Name())
		if isMappingFile(pathToFile) {
			result = append(result, pathToFile)
		}
	}

	return result, nil
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
			mappingFromStore.Tags = mappingFromFiles.Tags // Ensure new tags get propagated
			result[index] = mappingFromStore
		} else {
			result[index] = mappingFromFiles
		}

		validationErr := result[index].Validate()
		if validationErr != nil {
			return nil, fmt.Errorf("Error in mapping from file: %s\n%s", result[index].Origin, validationErr)
		}
	}

	return sortMappings(result), nil
}

func isGoProxyMappingFile(pathToFile string) bool {
	return isMappingFile(pathToFile) && (strings.HasSuffix(pathToFile, "go-proxy.yaml") || strings.HasSuffix(pathToFile, "go-proxy.yml"))
}

func isMappingFile(pathToFile string) bool {
	return filepath.Ext(pathToFile) == ".yml" || filepath.Ext(pathToFile) == ".yaml"
}

func loadAllMappings() ([]Mapping, error) {
	mappingFiles, findErr := findAllMappingsFiles()
	if findErr != nil {
		return nil, findErr
	}

	mappings := make([]Mapping, 0)
	for _, pathToFile := range mappingFiles {
		loadedMappings, loadingErr := loadMappingsFromFiles(pathToFile)
		if loadingErr != nil {
			return nil, loadingErr
		}
		mappings = append(mappings, loadedMappings...)
	}

	return mappings, nil
}
