package handler

import (
	"io"
	"net/http"

	"github.com/visola/go-proxy/pkg/upstream"
)

// HandleResult represents the result of handler handling a request
type HandleResult struct {
	Body         io.ReadCloser
	ErrorMessage string
	ExecutedURL  string
	Headers      map[string][]string
	ResponseCode int
}

type Handler interface {
	Handle(upstream.Mapping, http.Request) HandleResult
}

// Handlers contains all the available handlers mapped by the type of mapping
// that they can handle.
var Handlers = make(map[string]Handler)
