package upstream

// All possible mapping types
const (
	// HttpProxy maps an incoming request to an HTTP upstream
	HttpProxy = "HttpProxy"

	// StaticLocal maps an incoming request to a file in the local disk
	StaticLocal = "StaticLocal"
)

// Mapping represents a mapping route
type Mapping struct {
	From        string `json:"from"`
	Regexp      string `json:"regexp"`
	To          string `json:"to"`
	Type        string `json:"type"`
	UpsteamName string `json:"upstreamName"`
}
