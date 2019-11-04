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

func AddUpstreams(toAdd []Upstream) {
	upstreamsMutex.Lock()
	defer upstreamsMutex.Unlock()

	upstreams = append(upstreams, toAdd...)
}
