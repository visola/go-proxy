package upstream

// All possible mapping types
const (
	// MappingTypeHttpProxy maps an incoming request to an HTTP upstream
	MappingTypeHttpProxy = "HttpProxy"

	// MappingTypeStatic maps an incoming request to a file in the local disk
	MappingTypeStatic = "Static"
)

// Mapping represents a mapping route
type Mapping struct {
	From         string `json:"from"`
	Regexp       string `json:"regexp"`
	To           string `json:"to"`
	Type         string `json:"type"`
	UpstreamName string `json:"upstreamName"`
}
