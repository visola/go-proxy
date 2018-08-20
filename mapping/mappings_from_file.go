package mapping

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

type yamlMapping struct {
	Static []Mapping
	Proxy  []Mapping
}

func generateID(mapping Mapping) string {
	hasher := sha1.New()
	isProxy := "0"
	if mapping.Proxy {
		isProxy = "1"
	}
	headers := fmt.Sprintf("%s", mapping.Inject.Headers)
	hasher.Write([]byte(mapping.Origin + mapping.From + mapping.To + mapping.Regexp + isProxy + headers))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func loadMappingsFromFiles(mappingDir string, file os.FileInfo) ([]Mapping, error) {
	mappings := make([]Mapping, 0)

	mappingsFromFile, mappingErr := readMappingFile(path.Join(mappingDir, file.Name()))
	if mappingErr != nil {
		return nil, fmt.Errorf("Error while reading configuration file: %s\n%s", file.Name(), mappingErr)
	}

	for _, staticMapping := range mappingsFromFile.Static {
		staticMapping.Active = true
		staticMapping.Origin = file.Name()
		staticMapping.Proxy = false
		staticMapping.MappingID = generateID(staticMapping)

		mappings = append(mappings, staticMapping)
	}

	for _, proxyMapping := range mappingsFromFile.Proxy {
		proxyMapping.Active = true
		proxyMapping.Origin = file.Name()
		proxyMapping.Proxy = true
		proxyMapping.MappingID = generateID(proxyMapping)

		mappings = append(mappings, proxyMapping)
	}

	return mappings, nil
}

func readMappingFile(file string) (loadedMapping yamlMapping, err error) {
	var yamlContent []byte
	yamlContent, err = ioutil.ReadFile(file)

	if err != nil {
		return loadedMapping, err
	}

	err = yaml.Unmarshal(yamlContent, &loadedMapping)
	return loadedMapping, err
}
