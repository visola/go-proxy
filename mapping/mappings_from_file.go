package mapping

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type yamlMapping struct {
	Proxy  []Mapping
	Static []Mapping
	Tags   []string
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

	mappings = append(mappings, mappingFromYamlMapping(mappingsFromFile.Static, file.Name(), false, mappingsFromFile.Tags)...)
	mappings = append(mappings, mappingFromYamlMapping(mappingsFromFile.Proxy, file.Name(), true, mappingsFromFile.Tags)...)

	return mappings, nil
}

func mappingFromYamlMapping(fromYaml []Mapping, origin string, isProxy bool, extraTags []string) []Mapping {
	result := make([]Mapping, len(fromYaml))
	for i, m := range fromYaml {
		m.Active = true
		m.Origin = origin
		m.Proxy = isProxy
		m.MappingID = generateID(m)
		m.Tags = uniqueAndSort(append(m.Tags, extraTags...))
		m.Variables = m.GetVariables()

		result[i] = m
	}
	return result
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

func uniqueAndSort(arr []string) []string {
	if arr == nil {
		arr = make([]string, 0)
	}

	// Put everything in a map
	m := make(map[string]bool)
	for _, e := range arr {
		m[strings.ToLower(e)] = true
	}

	// Get all the keys back
	result := make([]string, len(m))
	count := 0
	for k := range m {
		result[count] = k
		count++
	}

	sort.Strings(result)
	return result
}
