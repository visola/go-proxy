package listener

import (
	"sync"
)

var (
	handledRequests      = make(map[string]HandleResult)
	handledRequestsMutex sync.Mutex
)

// RequestHandlingChanged stores or update the handle result
func RequestHandlingChanged(result HandleResult) {
	handledRequestsMutex.Lock()
	defer handledRequestsMutex.Unlock()

	handledRequests[result.ID] = result
}

// GetRequests return an array of handle result sorted by started time
func GetRequests() map[string]HandleResult {
	return handledRequests
}
