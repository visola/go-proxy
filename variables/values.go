package variables

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/Everbridge/go-proxy/configuration"
)

const valuesFile = ".values"

// GetPossibleValues returns all the possible values for all variables
func GetPossibleValues() (map[string][]string, error) {
	stateDir, stateDirErr := configuration.GetConfigurationDirectory()
	if stateDirErr != nil {
		return nil, stateDirErr
	}

	currentStateFile, currentStateFileErr := os.Open(path.Join(stateDir, valuesFile))
	if os.IsNotExist(currentStateFileErr) {
		return make(map[string][]string), nil
	}

	if currentStateFileErr != nil {
		return nil, currentStateFileErr
	}

	currentStateData, readErr := ioutil.ReadAll(currentStateFile)
	if readErr != nil {
		return nil, readErr
	}

	var result map[string][]string
	unmarshalErr := json.Unmarshal(currentStateData, &result)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return result, nil
}

// SetPossibleValues replace all possible values with the new set
func SetPossibleValues(toStore map[string][]string) error {
	stateDir, stateDirErr := configuration.GetConfigurationDirectory()
	if stateDirErr != nil {
		return stateDirErr
	}

	data, dataErr := json.Marshal(toStore)
	if dataErr != nil {
		return dataErr
	}

	return ioutil.WriteFile(path.Join(stateDir, valuesFile), data, 0644)
}
