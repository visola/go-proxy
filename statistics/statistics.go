package statistics

import (
	"github.com/visola/go-proxy/config"
)

// ProxiedRequest represents a request that was proxied
type ProxiedRequest struct {
	EndTime      int64          `json:"endTime"`
	Error        string         `json:"error"`
	ExecutedURL  string         `json:"executedURL"`
	Mapping      config.Mapping `json:"mapping"`
	RequestData  HTTPData       `json:"requestData"`
	Method       string         `json:"method"`
	ResponseCode int            `json:"responseCode"`
	ResponseData HTTPData       `json:"responseData"`
	RequestedURL string         `json:"requestedURL"`
	StartTime    int64          `json:"startTime"`
}

// HTTPData represents data sent or received in an HTTP call
type HTTPData struct {
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
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
