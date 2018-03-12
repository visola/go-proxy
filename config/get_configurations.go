package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

var configurations []Mapping

// GetConfigurations get all configurations or load them if not loaded so far
func GetConfigurations() ([]Mapping, error) {
	if configurations != nil {
		return configurations, nil
	}

	var err error
	configurations, err = loadAllConfigurations()
	return configurations, err
}

func getConfigDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.go-proxy", nil
}

func loadAllConfigurations() ([]Mapping, error) {
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

		loadedConfigurations, loadingErr := loadConfigurations(configDir, file)
		if loadingErr != nil {
			return nil, loadingErr
		}
		configurations = append(configurations, loadedConfigurations...)
	}

	sortConfigurations(configurations)
	return configurations, nil
}

func loadConfigurations(configDir string, file os.FileInfo) ([]Mapping, error) {
	configurations = make([]Mapping, 0)

	config, configErr := readConfiguration(path.Join(configDir, file.Name()))
	if configErr != nil {
		return nil, configErr
	}

	for _, staticConfig := range config.Static {
		configurations = append(configurations, fromYAMLMapping(staticConfig, file.Name(), false))
	}

	for _, staticConfig := range config.Proxy {
		configurations = append(configurations, fromYAMLMapping(staticConfig, file.Name(), true))
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

func sortConfigurations(configurations []Mapping) {
	sort.Slice(configurations, func(i, j int) bool {
		pathI := strings.ToLower(configurations[i].From)
		pathJ := strings.ToLower(configurations[j].From)

		if len(pathI) == len(pathJ) {
			return strings.Compare(pathI, pathJ) < 0
		}

		return len(pathI) > len(pathJ)
	})
}
