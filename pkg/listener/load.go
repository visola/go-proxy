package listener

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/visola/go-proxy/pkg/configuration"
	"gopkg.in/yaml.v2"
)

func loadFromFile(pathToFile string) (*Listener, error) {
	var err error
	var yamlContent []byte
	if yamlContent, err = ioutil.ReadFile(pathToFile); err != nil {
		return nil, err
	}

	var loadedListener Listener
	if err = yaml.Unmarshal(yamlContent, &loadedListener); err != nil {
		return nil, err
	}

	stats, statsErr := os.Stat(pathToFile)
	if statsErr != nil {
		return nil, statsErr
	}

	loadedListener.Origin = configuration.Origin{
		File:     pathToFile,
		LoadedAt: stats.ModTime().Unix(),
	}

	if loadedListener.Name == "" {
		fileName := filepath.Base(pathToFile)
		ext := path.Ext(fileName)
		loadedListener.Name = fileName[:len(fileName)-len(ext)]
	}

	return &loadedListener, nil
}
