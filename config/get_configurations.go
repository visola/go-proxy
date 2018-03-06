package config

import (
	"io/ioutil"
	"os/user"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

var configurations []Mapping

// GetConfigurations load all configurations.
func GetConfigurations() ([]Mapping, error) {
	if configurations != nil {
		return configurations, nil
	}

	configDir, configDirErr := getConfigDirectory()
	if configDirErr != nil {
		panic(configDirErr)
	}

	files, filesErr := ioutil.ReadDir(configDir)
	if filesErr != nil {
		panic(filesErr)
	}

	configurations = make([]Mapping, 0)

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".yml" && filepath.Ext(file.Name()) != ".yaml" {
			continue
		}

		config, configErr := readConfiguration(path.Join(configDir, file.Name()))
		if configErr != nil {
			return nil, configErr
		}

		for _, staticConfig := range config.Static {
			configurations = append(configurations, Mapping{
				From:   staticConfig.From,
				To:     staticConfig.To,
				Origin: file.Name(),
				Proxy:  false,
			})
		}

		for _, staticConfig := range config.Proxy {
			configurations = append(configurations, Mapping{
				From:   staticConfig.From,
				To:     staticConfig.To,
				Origin: file.Name(),
				Proxy:  true,
			})
		}
	}

	sort.Slice(configurations, func(i, j int) bool {
		pathI := strings.ToLower(configurations[i].From)
		pathJ := strings.ToLower(configurations[j].From)

		if len(pathI) == len(pathJ) {
			return strings.Compare(pathI, pathJ) < 0
		}

		return len(pathI) > len(pathJ)
	})
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

func getConfigDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.go-proxy", nil
}
