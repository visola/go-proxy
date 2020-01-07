package listener

import (
	"github.com/visola/go-proxy/pkg/event"
	"sync"
)

var (
	handledRequests      = make(map[string]event.HandleResult)
	handledRequestsMutex sync.Mutex
)

// RequestHandlingChanged stores or update the handle result
func RequestHandlingChanged(result event.HandleResult) {
	handledRequestsMutex.Lock()
	defer handledRequestsMutex.Unlock()

	handledRequests[result.ID] = result
	event.EmitHandleResult(result)
}

// GetRequests return an array of handle result sorted by started time
func GetRequests() map[string]event.HandleResult {
	return handledRequests
}
