package httputil

import (
	"net/http"
)

// NotFound sends an HTTP 404 response with a text/plain message
func NotFound(req *http.Request, w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(message))
}
