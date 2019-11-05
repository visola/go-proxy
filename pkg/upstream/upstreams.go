package upstream

import "sync"

// Upstream represents a place where requests can be directed to
type Upstream struct {
	Name   string
	Origin UpstreamOrigin
}

// UpstreamOrigin is where the upstream was loaded from
type UpstreamOrigin struct {
	File     string
	LoadedAt int64
}

var (
	upstreams      = make([]Upstream, 0)
	upstreamsMutex sync.Mutex
)

// AddUpstreams add an array of upstreams to the available array of upstreams
func AddUpstreams(toAdd []Upstream) {
	upstreamsMutex.Lock()
	defer upstreamsMutex.Unlock()

	upstreams = append(upstreams, toAdd...)
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
