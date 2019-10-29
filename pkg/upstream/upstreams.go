package upstream

// Upstream represents a place where requests can be directed to
type Upstream struct {
	Name   string
	Origin UpstreamOrigin
}

type UpstreamOrigin struct {
	File     string
	LoadedAt int64
}
