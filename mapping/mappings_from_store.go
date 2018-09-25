package mapping

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/Everbridge/go-proxy/configuration"
)

func getStoredState() (map[string]Mapping, error) {
	mappingDir, mappingDirErr := configuration.GetConfigurationDirectory()
	if mappingDirErr != nil {
		return nil, mappingDirErr
	}

	currentStateFile, currentStateFileErr := os.Open(path.Join(mappingDir, stateFileName))
	if os.IsNotExist(currentStateFileErr) {
		return make(map[string]Mapping, 0), nil
	}

	if currentStateFileErr != nil {
		return nil, currentStateFileErr
	}

	currentStateData, readErr := ioutil.ReadAll(currentStateFile)
	if readErr != nil {
		return nil, readErr
	}

	var result map[string]Mapping
	unmarshalErr := json.Unmarshal(currentStateData, &result)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return result, nil
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

	mappingDir, mappingDirErr := configuration.GetConfigurationDirectory()
	if mappingDirErr != nil {
		return mappingDirErr
	}

	return ioutil.WriteFile(path.Join(mappingDir, stateFileName), data, 0644)
}
