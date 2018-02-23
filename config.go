package main

import (
	"io/ioutil"
	"os/user"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type proxyConfig struct {
	Static []mapping
	Proxy  []mapping
}

type mapping struct {
	From string
	To   string
}

type configMapping struct {
	from   string
	origin string
	proxy  bool
	to     string
}

func getConfigDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.go-proxy", nil
}

func getConfigurations() ([]configMapping, error) {
	configDir, configDirErr := getConfigDirectory()
	if configDirErr != nil {
		panic(configDirErr)
	}

	files, filesErr := ioutil.ReadDir(configDir)
	if filesErr != nil {
		panic(filesErr)
	}

	var configurations = make([]configMapping, 0)

	for _, file := range files {
		config, configErr := readConfiguration(path.Join(configDir, file.Name()))
		if configErr != nil {
			return nil, configErr
		}

		for _, staticConfig := range config.Static {
			configurations = append(configurations, configMapping{
				from:   staticConfig.From,
				to:     staticConfig.To,
				origin: file.Name(),
				proxy:  false,
			})
		}

		for _, staticConfig := range config.Proxy {
			configurations = append(configurations, configMapping{
				from:   staticConfig.From,
				to:     staticConfig.To,
				origin: file.Name(),
				proxy:  true,
			})
		}
	}

	sort.Slice(configurations, func(i, j int) bool {
		pathI := strings.ToLower(configurations[i].from)
		pathJ := strings.ToLower(configurations[j].from)

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
