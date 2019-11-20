package handler

import (
	"net/http"

	"github.com/visola/go-proxy/pkg/upstream"
)

// HandleResult represents the result of handler handling a request
type HandleResult struct {
	Body         string              `json:"body"`
	ErrorMessage string              `json:"errorMessage"`
	ExecutedURL  string              `json:"executedUrl"`
	Headers      map[string][]string `json:"headers"`
	ResponseCode int                 `json:"responseCode"`
}

type Handler interface {
	Handle(upstream.Mapping, http.Request) HandleResult
	Matches(upstream.Mapping, http.Request) bool
}

// Handlers contains all the available handlers mapped by the type of mapping
// that they can handle.
var Handlers = make(map[string]*Handler)
