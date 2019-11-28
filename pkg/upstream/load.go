package upstream

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type upstreamFile struct {
	Static    []yamlMapping
	Upstreams []yamlUpstream
}

type yamlUpstream struct {
	Name   string
	Static []yamlMapping
}

type yamlMapping struct {
	From   string
	Regexp string
	To     string
}

func (m yamlMapping) toMapping(mappingType string, upstreamName string) Mapping {
	return Mapping{
		From:         m.From,
		Regexp:       m.Regexp,
		Type:         mappingType,
		UpstreamName: upstreamName,
	}
}

func loadFromFile(pathToFile string) (upstreams []Upstream, err error) {
	var yamlContent []byte
	if yamlContent, err = ioutil.ReadFile(pathToFile); err != nil {
		return
	}

	var yamlFile upstreamFile
	if err = yaml.Unmarshal(yamlContent, &yamlFile); err != nil {
		return
	}

	stats, statsErr := os.Stat(pathToFile)

	if statsErr != nil {
		return upstreams, statsErr
	}

	origin := UpstreamOrigin{
		File:     pathToFile,
		LoadedAt: stats.ModTime().Unix(),
	}

	upstreams = make([]Upstream, 0)

	rootUpstream := Upstream{
		Mappings: make([]Mapping, 0),
		Name:     nameFromFilePath(pathToFile),
		Origin:   origin,
	}

	for _, m := range yamlFile.Static {
		rootUpstream.Mappings = append(rootUpstream.Mappings, m.toMapping(MappingTypeStatic, rootUpstream.Name))
	}

	upstreams = append(upstreams, rootUpstream)

	for _, u := range yamlFile.Upstreams {
		innerUpstream := Upstream{
			Mappings: make([]Mapping, 0),
			Name:     u.Name,
			Origin:   origin,
		}

		for _, m := range u.Static {
			innerUpstream.Mappings = append(innerUpstream.Mappings, m.toMapping(MappingTypeStatic, innerUpstream.Name))
		}

		upstreams = append(upstreams, innerUpstream)
	}

	return
}

func nameFromFilePath(filePath string) string {
	ext := path.Ext(filePath)
	filePath = filePath[:len(filePath)-len(ext)]

	baseFile := path.Base(filePath)
	if baseFile != "go-proxy" {
		return baseFile
	}

	parent := filePath[:len(filePath)-len(baseFile)]
	return path.Base(parent)
}
