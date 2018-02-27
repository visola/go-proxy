package config

// Mapping represents a mapping configuration loaded from some file.
type Mapping struct {
	From   string `json:"from"`
	Origin string `json:"origin"`
	Proxy  bool   `json:"proxy"`
	To     string `json:"to"`
}
