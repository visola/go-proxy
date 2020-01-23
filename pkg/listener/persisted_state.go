package listener

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/visola/go-proxy/pkg/configuration"
)

const listenerSubDirectory = "listeners"

// LoadFromPersistedState loads all listeners from the files in the configuration directory into memory
func LoadFromPersistedState() error {
	configDir, configDirErr := configuration.GetConfigurationDirectory()
	if configDirErr != nil {
		return configDirErr
	}

	listenersDir := filepath.Join(configDir, listenerSubDirectory)
	if _, err := os.Stat(listenersDir); os.IsNotExist(err) {
		return nil
	}

	files, readDirErr := ioutil.ReadDir(listenersDir)
	if readDirErr != nil {
		return readDirErr
	}

	loadedListeners := make(map[string]*Listener)

	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == ".yml" || ext == ".yaml" {
				loadedListener, loadErr := loadFromFile(filepath.Join(listenersDir, file.Name()))
				if loadErr != nil {
					return loadErr
				}

				loadedListeners[loadedListener.Name] = loadedListener
			}
		}
	}

	currentListeners = loadedListeners
	return nil
}

// SaveToPersistedState save all listeners back to files
func SaveToPersistedState() error {
	currentListenersMutex.Lock()
	defer currentListenersMutex.Unlock()

	for _, l := range currentListeners {
		if saveErr := Save(l); saveErr != nil {
			return saveErr
		}
	}

	return nil
}
