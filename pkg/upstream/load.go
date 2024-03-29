package upstream

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/visola/go-proxy/pkg/configuration"
	"gopkg.in/yaml.v2"
)

type upstreamFile struct {
	Proxy     []yamlProxyEndpoint
	Static    []yamlStaticEndpoint
	Upstreams map[string]yamlUpstream
}

type yamlUpstream struct {
	Proxy  []yamlProxyEndpoint
	Static []yamlStaticEndpoint
}

type yamlProxyEndpoint struct {
	From    string
	Headers map[string]arrayOrString
	Regexp  string
	To      string
}

func (m yamlProxyEndpoint) toEndpoint(upstreamName string) *ProxyEndpoint {
	return &ProxyEndpoint{
		Headers: FromMapOfArrayOfStrings(m.Headers),
		To:      m.To,
		BaseEndpoint: BaseEndpoint{
			From:         m.From,
			Regexp:       m.Regexp,
			UpstreamName: upstreamName,
		},
	}
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

	origin := configuration.Origin{
		File:     pathToFile,
		LoadedAt: stats.ModTime().Unix(),
	}

	upstreams = make([]Upstream, 0)

	rootUpstream := Upstream{
		ProxyEndpoints:  make([]*ProxyEndpoint, 0),
		StaticEndpoints: make([]*StaticEndpoint, 0),
		Name:            nameFromFilePath(pathToFile),
		Origin:          origin,
	}

	for _, m := range yamlFile.Proxy {
		rootUpstream.ProxyEndpoints = append(rootUpstream.ProxyEndpoints, m.toEndpoint(rootUpstream.Name))
	}

	for _, m := range yamlFile.Static {
		rootUpstream.StaticEndpoints = append(rootUpstream.StaticEndpoints, m.toEndpoint(rootUpstream.Name))
	}

	upstreams = append(upstreams, rootUpstream)

	for name, u := range yamlFile.Upstreams {
		innerUpstream := Upstream{
			ProxyEndpoints:  make([]*ProxyEndpoint, 0),
			StaticEndpoints: make([]*StaticEndpoint, 0),
			Name:            name,
			Origin:          origin,
		}

		for _, m := range u.Proxy {
			innerUpstream.ProxyEndpoints = append(innerUpstream.ProxyEndpoints, m.toEndpoint(innerUpstream.Name))
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
