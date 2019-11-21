package upstream

import "sync"

// Upstream represents a place where requests can be directed to
type Upstream struct {
	Mappings []Mapping      `json:"mappings"`
	Name     string         `json:"name"`
	Origin   UpstreamOrigin `json:"origin"`
}

// UpstreamOrigin is where the upstream was loaded from
type UpstreamOrigin struct {
	File     string `json:"file"`
	LoadedAt int64  `json:"loadedAt"`
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
		upstreams[toAdd.Name] = toAdd
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
