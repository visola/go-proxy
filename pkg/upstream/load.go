package upstream

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type upstreamFile struct {
	Upstreams []yamlUpstream
}

type yamlUpstream struct {
	Name string
}

func LoadFromFile(pathToFile string) (upstreams []Upstream, err error) {
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

	upstreams = append(upstreams, Upstream{
		Name:   nameFromFilePath(pathToFile),
		Origin: origin,
	})

	for _, u := range yamlFile.Upstreams {
		upstreams = append(upstreams, Upstream{
			Name:   u.Name,
			Origin: origin,
		})
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
