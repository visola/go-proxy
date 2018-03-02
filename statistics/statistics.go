package statistics

import (
	"github.com/visola/go-proxy/config"
)

// ProxiedRequest represents a request that was proxied
type ProxiedRequest struct {
	EndTime       int64          `json:"endTime"`
	Error         string         `json:"error"`
	Mapping       config.Mapping `json:"mapping"`
	ProxiedTo     string         `json:"proxiedTo"`
	ResponseCode  int            `json:"responseCode"`
	RequestedPath string         `json:"requestedPath"`
	StartTime     int64          `json:"startTime"`
}

const maxRequestsToKeep = 1000

var (
	channel         = make(chan ProxiedRequest)
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

func addToArray() {
	if reading == true {
		return
	}

	reading = true
	for {
		proxiedRequests = append(proxiedRequests, <-channel)
		if len(proxiedRequests) > maxRequestsToKeep {
			proxiedRequests = proxiedRequests[1:]
		}
	}
}
