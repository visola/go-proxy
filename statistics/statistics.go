package statistics

import (
	"github.com/Everbridge/go-proxy/mapping"
)

// ProxiedRequest represents a request that was proxied
type ProxiedRequest struct {
	EndTime      int64           `json:"endTime,omitempty"`
	Error        string          `json:"error,omitempty"`
	ExecutedURL  string          `json:"executedURL,omitempty"`
	Mapping      mapping.Mapping `json:"mapping,omitempty"`
	RequestData  HTTPData        `json:"requestData,omitempty"`
	Method       string          `json:"method,omitempty"`
	ResponseCode int             `json:"responseCode,omitempty"`
	ResponseData HTTPData        `json:"responseData,omitempty"`
	RequestedURL string          `json:"requestedURL,omitempty"`
	StartTime    int64           `json:"startTime,omitempty"`
}

// HTTPData represents data sent or received in an HTTP call
type HTTPData struct {
	Headers map[string][]string `json:"headers,omitempty"`
	Body    string              `json:"body,omitempty"`
}

const maxRequestsToKeep = 1000

var (
	channel         = make(chan ProxiedRequest)
	listeners       = make([]func(ProxiedRequest), 0)
	proxiedRequests = make([]ProxiedRequest, 0)
	reading         = false
)

// AddProxiedRequest add the proxied request to the channel
func AddProxiedRequest(proxiedRequest ProxiedRequest) {
	go addToArray()
	channel <- proxiedRequest
}

// GetProxiedRequests return all the proxied requests so far
func GetProxiedRequests() []ProxiedRequest {
	return proxiedRequests
}

// OnRequestProxied will register the callback to be called when a new request
// is proxied.
func OnRequestProxied(callback func(ProxiedRequest)) {
	listeners = append(listeners, callback)
}

func addToArray() {
	if reading == true {
		return
	}

	reading = true
	for {
		proxiedRequest := <-channel
		proxiedRequests = append(proxiedRequests, proxiedRequest)

		for _, listener := range listeners {
			listener(proxiedRequest)
		}

		if len(proxiedRequests) > maxRequestsToKeep {
			proxiedRequests = proxiedRequests[1:]
		}
	}
}
