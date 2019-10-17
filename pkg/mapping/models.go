package mapping

// Possible types of mappings
const (
	MappingTypeProxy  = "PROXY"
	MappingTypeStatic = "STATIC"
)

// Injection represents parameters that can be injected into proxied requests
type Injection struct {
	Headers map[string]string `json:"headers"`
}

// Mapping represents a mapping
type Mapping struct {
	File     string    `json:"file"`
	From     string    `json:"from"`
	LoadedAt int64     `json:"loadedAt"`
	Inject   Injection `json:"injection"`
	Regexp   string    `json:"regexp"`
	To       string    `json:"to"`
	Type     string    `json:"type"`
}

func (m *Mapping) Initialize(filePath string, mappingType string, loadedAt int64) {
	m.File = filePath
	m.LoadedAt = loadedAt
	m.Type = mappingType
}
