package configuration

import (
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const configFileName = ".config"

// LoadConfiguration loads the current user configuration
func LoadConfiguration() (*Configuration, error) {
	configDir, dirErr := GetConfigurationDirectory()
	if dirErr != nil {
		return nil, dirErr
	}

	result := new(Configuration)
	result.BaseDirectories = make([]string, 0)

	configContent, readErr := ioutil.ReadFile(path.Join(configDir, configFileName))
	if os.IsNotExist(readErr) {
		return result, nil
	}

	if readErr != nil {
		return nil, readErr
	}

	unmarshalErr := yaml.Unmarshal(configContent, result)
	return result, unmarshalErr
}

// SaveConfiguration saves a configuration to the file
func SaveConfiguration(toSave Configuration) error {
	data, marshalErr := yaml.Marshal(toSave)
	if marshalErr != nil {
		return marshalErr
	}

	configDir, dirErr := GetConfigurationDirectory()
	if dirErr != nil {
		return dirErr
	}

	return ioutil.WriteFile(path.Join(configDir, configFileName), data, 0644)
}
