package variables

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/visola/go-proxy/configuration"
)

const valuesFile = ".values"

// ValuesState is the struct used to store the current values state in the file
type ValuesState struct {
	PossibleValues map[string][]string `json:"possibleValues"`
	SelectedValues map[string]string   `json:"selectedValues"`
}

// GetPossibleValues returns all the possible values for all variables
func GetPossibleValues() (map[string][]string, error) {
	currentState, err := loadCurrentState()
	if err != nil {
		return nil, err
	}

	return currentState.PossibleValues, nil
}

// GetSelectedValues returns all selected values for each variable
func GetSelectedValues() (map[string]string, error) {
	currentState, err := loadCurrentState()
	if err != nil {
		return nil, err
	}

	return currentState.SelectedValues, nil
}

// SetPossibleValues replace all possible values with the new set
func SetPossibleValues(toStore map[string][]string) error {
	currentState, err := loadCurrentState()
	if err != nil {
		return err
	}

	currentState.PossibleValues = toStore
	return saveCurrentState(currentState)
}

// SetValue sets a value for a specific variable
func SetValue(varName string, value string) error {
	currentState, err := loadCurrentState()
	if err != nil {
		return err
	}

	currentState.SelectedValues[varName] = value

	// Add value to possible values if not present
	present := false
	for _, pV := range currentState.PossibleValues[varName] {
		if pV == value {
			present = true
			break
		}
	}

	if !present {
		currentState.PossibleValues[varName] = append(currentState.PossibleValues[varName], value)
	}

	return saveCurrentState(currentState)
}

func loadCurrentState() (*ValuesState, error) {
	stateDir, stateDirErr := configuration.GetConfigurationDirectory()
	if stateDirErr != nil {
		return nil, stateDirErr
	}

	currentStateFile, currentStateFileErr := os.Open(path.Join(stateDir, valuesFile))
	if os.IsNotExist(currentStateFileErr) {
		return &ValuesState{
			PossibleValues: make(map[string][]string),
			SelectedValues: make(map[string]string),
		}, nil
	}

	if currentStateFileErr != nil {
		return nil, currentStateFileErr
	}

	currentStateData, readErr := ioutil.ReadAll(currentStateFile)
	if readErr != nil {
		return nil, readErr
	}

	var result ValuesState
	unmarshalErr := json.Unmarshal(currentStateData, &result)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &result, nil
}

func saveCurrentState(currentState *ValuesState) error {
	stateDir, stateDirErr := configuration.GetConfigurationDirectory()
	if stateDirErr != nil {
		return stateDirErr
	}

	data, dataErr := json.Marshal(currentState)
	if dataErr != nil {
		return dataErr
	}

	return ioutil.WriteFile(path.Join(stateDir, valuesFile), data, 0644)
}
