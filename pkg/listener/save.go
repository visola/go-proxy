package listener

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/visola/go-proxy/pkg/configuration"
	"gopkg.in/yaml.v2"
)

// Save saves a listener to it's file, or to the default file if not specified
func Save(toSave *Listener) error {
	pathToSave := toSave.Origin.File
	if pathToSave == "" {
		defaultPath, err := getDefaultPath()
		if err != nil {
			return err
		}
		pathToSave = defaultPath
	}

	pathDir := filepath.Dir(pathToSave)
	if mkdirErr := os.MkdirAll(pathDir, 0766); mkdirErr != nil {
		return mkdirErr
	}

	data, marshalErr := yaml.Marshal(toSave)
	if marshalErr != nil {
		return marshalErr
	}

	toSave.Origin = configuration.Origin{
		File:     pathToSave,
		LoadedAt: time.Now().Unix(),
	}

	currentListeners[toSave.Name] = toSave
	return ioutil.WriteFile(pathToSave, data, 0644)
}

func getDefaultPath() (string, error) {
	configDir, configDirErr := configuration.GetConfigurationDirectory()
	if configDirErr != nil {
		return "", configDirErr
	}

	return filepath.Join(configDir, listenerSubDirectory, defaultFile), nil
}
