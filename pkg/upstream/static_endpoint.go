package upstream

import (
	"net/http"
)

// StaticEndpoint is an endpoint that responds with file from disk
type StaticEndpoint struct {
	BaseEndpoint
	To string `json:"to"`
}

func (s StaticEndpoint) Handle(http.Request) HandleResult {
	// TODO - Implement this handling method
	return HandleResult{
		ResponseCode: http.StatusOK,
	}
}
