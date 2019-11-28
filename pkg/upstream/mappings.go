package upstream

import (
	"net/http"
	"regexp"
	"strings"
)

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
	Type         string `json:"type"`
	UpstreamName string `json:"upstreamName"`
}

// Matches check if the request matches the request
func (m Mapping) Matches(req http.Request) bool {
	if m.From != "" && strings.HasPrefix(req.URL.Path, m.From) {
		return true
	}

	if m.Regexp != "" {
		r, err := regexp.Compile(m.Regexp)
		if err != nil {
			return false
		}

		return r.MatchString(req.URL.Path)
	}

	return false
}
