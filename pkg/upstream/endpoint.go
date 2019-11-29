package upstream

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

// HandleResult represents the result of an endpoint handling a request
type HandleResult struct {
	Body         io.ReadCloser
	ErrorMessage string
	ExecutedURL  string
	Headers      map[string][]string
	ResponseCode int
}

// Endpoint represents one route inside an Upstream
type Endpoint interface {
	Handle(http.Request) HandleResult
	Matches(http.Request) bool
}

// BaseEndpoint represents a base endpoint route
type BaseEndpoint struct {
	From         string `json:"from"`
	Regexp       string `json:"regexp"`
	UpstreamName string `json:"upstreamName"`
}

// Matches check if the request matches the request
func (m BaseEndpoint) Matches(req http.Request) bool {
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
