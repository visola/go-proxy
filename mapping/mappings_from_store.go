package mapping

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

func getStoredState() (map[string]Mapping, error) {
	mappingDir, mappingDirErr := getMappingDirectory()
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
