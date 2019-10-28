package loader

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/visola/go-proxy/pkg/upstream"
	"gopkg.in/yaml.v2"
)

type upstreamFile struct {
	Upstreams []yamlUpstream
}

type yamlUpstream struct {
	Name string
}

func loadUpstreamsFromFile(pathToFile string) (upstreams []upstream.Upstream, err error) {
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

	origin := upstream.UpstreamOrigin{
		File:     pathToFile,
		LoadedAt: stats.ModTime().Unix(),
	}

	upstreams = make([]upstream.Upstream, 0)

	upstreams = append(upstreams, upstream.Upstream{
		Name:   nameFromFilePath(pathToFile),
		Origin: origin,
	})

	for _, u := range yamlFile.Upstreams {
		upstreams = append(upstreams, upstream.Upstream{
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
