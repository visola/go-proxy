package main

import (
	"io/ioutil"
	"os/user"
	"path"

	"gopkg.in/yaml.v2"
)

type proxyConfig struct {
	Static []staticConfiguration
}

type staticConfiguration struct {
	From string
	To   string
}

func getConfigDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.go-proxy", nil
}

func getConfigurations() ([]proxyConfig, error) {
	configDir, configDirErr := getConfigDirectory()
	if configDirErr != nil {
		panic(configDirErr)
	}

	files, filesErr := ioutil.ReadDir(configDir)
	if filesErr != nil {
		panic(filesErr)
	}

	var configurations = make([]proxyConfig, 0)
	for _, file := range files {
		config, configErr := readConfiguration(path.Join(configDir, file.Name()))
		if configErr != nil {
			return nil, configErr
		}
		configurations = append(configurations, config)
	}

	return configurations, nil
}

func readConfiguration(file string) (loadedConfig proxyConfig, err error) {
	var yamlContent []byte
	yamlContent, err = ioutil.ReadFile(file)

	if err != nil {
		return loadedConfig, err
	}

	err = yaml.Unmarshal(yamlContent, &loadedConfig)
	return loadedConfig, err
}
