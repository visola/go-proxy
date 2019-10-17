package mapping

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type yamlDocument struct {
	Proxy  []Mapping
	Static []Mapping
}

func loadMappingFromFile(pathToFile string) (mappings []Mapping, err error) {
	var yamlContent []byte
	yamlContent, err = ioutil.ReadFile(pathToFile)

	if err != nil {
		return mappings, err
	}

	mappings = make([]Mapping, 0)

	var doc yamlDocument
	err = yaml.Unmarshal(yamlContent, &doc)

	loadedAt := time.Now().Unix()

	for _, proxyMapping := range doc.Proxy {
		proxyMapping.Initialize(pathToFile, MappingTypeProxy, loadedAt)
		mappings = append(mappings, proxyMapping)
	}

	for _, staticMapping := range doc.Static {
		staticMapping.Initialize(pathToFile, MappingTypeStatic, loadedAt)
		mappings = append(mappings, staticMapping)
	}

	return mappings, err
}
