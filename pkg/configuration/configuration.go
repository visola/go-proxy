package configuration

import (
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/visola/go-proxy/pkg/listener"
	"gopkg.in/yaml.v2"
)

// Configuration represents a configuration state for the proxy
type Configuration struct {
	Listeners map[int]listener.Listener
}

const currentStateFile = ".current_configuration"

var currentStateFileMutex sync.Mutex

// LoadFromPersistedState initializes the current configuration from the persisted state
func LoadFromPersistedState() error {
	configDir, configDirErr := GetConfigurationDirectory()
	if configDirErr != nil {
		return configDirErr
	}

	config, configErr := loadConfiguration(path.Join(configDir, currentStateFile))
	if configErr != nil {
		return configErr
	}

	currentConfiguration := &config
	if currentConfiguration.Listeners == nil {
		currentConfiguration.Listeners = make(map[int]listener.Listener)
	}

	// Initialize all listeners as inactive
	for _, l := range currentConfiguration.Listeners {
		l.Active = false
	}

	listener.SetListeners(currentConfiguration.Listeners)
	return nil
}

// SaveToPersistedState saves the current state to the persisted state file in the configuration directory
func SaveToPersistedState() error {
	configDir, configDirErr := GetConfigurationDirectory()
	if configDirErr != nil {
		return configDirErr
	}

	toSave := Configuration{
		Listeners: listener.Listeners(),
	}

	currentStateFileMutex.Lock()
	defer currentStateFileMutex.Unlock()

	return saveConfiguration(toSave, path.Join(configDir, currentStateFile))
}

func loadConfiguration(statePath string) (Configuration, error) {
	result := Configuration{}

	configContent, readErr := ioutil.ReadFile(statePath)
	if os.IsNotExist(readErr) {
		return result, nil
	}

	if readErr != nil {
		return result, readErr
	}

	unmarshalErr := yaml.Unmarshal(configContent, &result)
	return result, unmarshalErr
}

func saveConfiguration(toSave Configuration, statePath string) error {
	data, marshalErr := yaml.Marshal(toSave)
	if marshalErr != nil {
		return marshalErr
	}

	return ioutil.WriteFile(statePath, data, 0644)
}
