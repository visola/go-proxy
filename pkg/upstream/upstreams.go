package upstream

import (
	"sync"

	"github.com/visola/go-proxy/pkg/configuration"
)

// Upstream represents a place where requests can be directed to
type Upstream struct {
	Name            string               `json:"name"`
	Origin          configuration.Origin `json:"origin"`
	ProxyEndpoints  []*ProxyEndpoint     `json:"proxyEndpoints"`
	StaticEndpoints []*StaticEndpoint    `json:"staticEndpoints"`
}

// Endpoints returns all the endpoints available for the upstream
func (u Upstream) Endpoints() []Endpoint {
	allEndpoints := make([]Endpoint, 0)

	for _, m := range u.ProxyEndpoints {
		allEndpoints = append(allEndpoints, m)
	}

	for _, m := range u.StaticEndpoints {
		allEndpoints = append(allEndpoints, m)
	}

	return allEndpoints
}

// Stores upstreams in a map using the name as key
var (
	upstreams      = make(map[string]Upstream, 0)
	upstreamsMutex sync.Mutex
)

// AddUpstreams add an array of upstreams to the available upstreams
func AddUpstreams(allToAdd []Upstream) {
	upstreamsMutex.Lock()
	defer upstreamsMutex.Unlock()

	for _, toAdd := range allToAdd {
		_, exists := upstreams[toAdd.Name]
		if !exists {
			upstreams[toAdd.Name] = toAdd
		}
	}
}

// Upstreams return a map containing all upstreams loaded
func Upstreams() map[string]Upstream {
	return upstreams
}

// UpstreamsPerFile returns a map with all upstreams grouped by file where they were loaded from
func UpstreamsPerFile() map[string][]Upstream {
	upstreamsPerFile := make(map[string][]Upstream)

	for _, oldUpstream := range upstreams {
		upstreamsInFile, exists := upstreamsPerFile[oldUpstream.Origin.File]
		if !exists {
			upstreamsInFile = make([]Upstream, 0)
		}
		upstreamsInFile = append(upstreamsInFile, oldUpstream)
		upstreamsPerFile[oldUpstream.Origin.File] = upstreamsInFile
	}

	return upstreamsPerFile
}
