package event

import (
	"sync"
)

// All event types will require an array of channels
var (
	requestListeners       = make([]HandleResultConsumer, 0)
	requestListenersMutext sync.Mutex
)

// HandleResultConsumer consumes handle result messages
type HandleResultConsumer interface {
	Consume(HandleResult)
}

// HandleBodyAndHeaders stores headers and body for a request or response
type HandleBodyAndHeaders struct {
	Body    []byte              `json:"body"`
	Headers map[string][]string `json:"headers"`
}

// HandleResult stores information about handled requests
type HandleResult struct {
	Error       string `json:"error"`
	ExecutedURL string `json:"executedURL"`
	ID          string `json:"id"`
	HandledBy   string `json:"handledBy"`
	StatusCode  int    `json:"statusCode"`
	URL         string `json:"url"`

	Request  HandleBodyAndHeaders `json:"request"`
	Response HandleBodyAndHeaders `json:"response"`

	Timings HandleTimings `json:"timings"`
}

// HandleTimings stores information about the timings
type HandleTimings struct {
	Completed int64 `json:"completed"`
	Handled   int64 `json:"handled"`
	Matched   int64 `json:"matched"`
	Started   int64 `json:"started"`
}

// AddRequestListener adds a listener to handle request results
func AddRequestListener(l HandleResultConsumer) {
	requestListenersMutext.Lock()
	defer requestListenersMutext.Unlock()

	requestListeners = append(requestListeners, l)
}

// EmitHandleResult emits a handle result event to all listeners
func EmitHandleResult(h HandleResult) {
	for _, l := range requestListeners {
		l.Consume(h)
	}
}

// RemoveRequestListener removes a listener from the list
func RemoveRequestListener(toRemove HandleResultConsumer) {
	requestListenersMutext.Lock()
	defer requestListenersMutext.Unlock()

	newListeners := make([]HandleResultConsumer, len(requestListeners)-1)
	counter := 0
	for _, l := range requestListeners {
		if l == toRemove {
			continue
		}

		newListeners[counter] = l
		counter++
	}

	requestListeners = newListeners
}
