package upstream

import "net/http"

// HandleResult represents the result of handler handling a request
type HandleResult struct {
	Body         string              `json:"body"`
	ErrorMessage string              `json:"errorMessage"`
	ExecutedURL  string              `json:"executedUrl"`
	Headers      map[string][]string `json:"headers"`
	ResponseCode int                 `json:"responseCode"`
}

type Handler interface {
	handle(Mapping, http.Request) HandleResult
	matches(Mapping, http.Request) bool
}
