package configuration

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/visola/go-proxy/pkg/listener"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Listeners map[int]listener.Listener
}

const currentStateFile = ".current_configuration"

// LoadFromPersistedState initializes the current configuration from the persisted state
func LoadFromPersistedState() error {
	configDir, configDirErr := GetConfigurationDirectory()
	if configDirErr != nil {
		return configDirErr
	}

	config, configErr := loadConfiguration(configDir, currentStateFile)
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

func loadConfiguration(configDir string, statePath string) (Configuration, error) {
	result := Configuration{}

	configContent, readErr := ioutil.ReadFile(path.Join(configDir, currentStateFile))
	if os.IsNotExist(readErr) {
		return result, nil
	}

	if readErr != nil {
		return result, readErr
	}

	unmarshalErr := yaml.Unmarshal(configContent, &result)
	return result, unmarshalErr
}
