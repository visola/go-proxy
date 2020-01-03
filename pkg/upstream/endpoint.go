package upstream

import (
	"net/http"
	"regexp"
	"strings"
)

// HandleResult represents the result of an endpoint handling a request
type HandleResult struct {
	ErrorMessage    string
	ExecutedURL     string
	ResponseBody    []byte
	ResponseHeaders map[string][]string
	ResponseCode    int
}

// Endpoint represents one route inside an Upstream
type Endpoint interface {
	Handle(*http.Request, http.ResponseWriter) HandleResult
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
