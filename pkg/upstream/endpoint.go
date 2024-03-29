package upstream

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

// Endpoint represents one route inside an Upstream
type Endpoint interface {
	Handle(*http.Request) (int, string, map[string][]string, io.ReadCloser)
	Matches(*http.Request) bool
	Path() string
}

// BaseEndpoint represents a base endpoint route
type BaseEndpoint struct {
	From         string `json:"from"`
	Regexp       string `json:"regexp"`
	UpstreamName string `json:"upstreamName"`

	compiledRegexp *regexp.Regexp
}

// Matches check if the request matches the request
func (m *BaseEndpoint) Matches(req *http.Request) bool {
	if m.From != "" && strings.HasPrefix(req.URL.Path, m.From) {
		return true
	}

	if r := m.ensureRegexp(); r != nil {
		return r.MatchString(req.URL.Path)
	}

	return false
}

// Path returns the matching path for an endpoint
func (m *BaseEndpoint) Path() string {
	path := m.From
	if path == "" {
		path = m.Regexp
	}
	return path
}

func (m *BaseEndpoint) ensureRegexp() *regexp.Regexp {
	if m.Regexp != "" {
		r, err := regexp.Compile(m.Regexp)
		if err != nil {
			return nil
		}

		m.compiledRegexp = r
	}

	return m.compiledRegexp
}
