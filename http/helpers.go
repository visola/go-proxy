package http

import (
	"fmt"
	"net/http"
)

// NotFound sends an HTTP 404 response with a text/plain message
func NotFound(req *http.Request, w http.ResponseWriter, path string) {
	if path == "" {
		path = req.URL.String()
	}

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Sorry, nothing here for: '%s'\n", path)))
}

// InternalError sends an HTTP 500 with a text/plain message showing the error message
// It also logs the message.
func InternalError(req *http.Request, w http.ResponseWriter, err error) {
	fmt.Printf("Error: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Sorry, something went wrong: '%s'\n", err)))
}
