package configuration

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/visola/go-proxy/pkg/listener"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Listeners map[int]listener.Listener
}

var (
	currentConfiguration      = Configuration{}
	currentConfigurationMutex sync.Mutex
)

const currentStateFile = ".current_configuration"

// Initialize initializes the current configuration from the persisted state
func Initialize() {
	if err := loadFromPersistedState(); err != nil {
		log.Fatalf("Error while loading persisted configuration: %v", err)
	}
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

func loadFromPersistedState() error {
	configDir, configDirErr := GetConfigurationDirectory()
	if configDirErr != nil {
		return configDirErr
	}

	currentConfigurationMutex.Lock()
	defer currentConfigurationMutex.Unlock()

	config, configErr := loadConfiguration(configDir, currentStateFile)
	if configErr != nil {
		return configErr
	}

	currentConfiguration = config
	return nil
}
