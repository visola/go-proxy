package upstream

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type upstreamFile struct {
	Static    []yamlStaticEndpoint
	Upstreams []yamlUpstream
}

type yamlUpstream struct {
	Name   string
	Static []yamlStaticEndpoint
}

type yamlStaticEndpoint struct {
	From   string
	Regexp string
	To     string
}

func (m yamlStaticEndpoint) toEndpoint(upstreamName string) *StaticEndpoint {
	return &StaticEndpoint{
		To: m.To,
		BaseEndpoint: BaseEndpoint{
			From:         m.From,
			Regexp:       m.Regexp,
			UpstreamName: upstreamName,
		},
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
		StaticEndpoints: make([]*StaticEndpoint, 0),
		Name:            nameFromFilePath(pathToFile),
		Origin:          origin,
	}

	for _, m := range yamlFile.Static {
		rootUpstream.StaticEndpoints = append(rootUpstream.StaticEndpoints, m.toEndpoint(rootUpstream.Name))
	}

	upstreams = append(upstreams, rootUpstream)

	for _, u := range yamlFile.Upstreams {
		innerUpstream := Upstream{
			StaticEndpoints: make([]*StaticEndpoint, 0),
			Name:            u.Name,
			Origin:          origin,
		}

		for _, m := range u.Static {
			innerUpstream.StaticEndpoints = append(innerUpstream.StaticEndpoints, m.toEndpoint(innerUpstream.Name))
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
