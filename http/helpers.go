package http

import (
	"fmt"
	"net/http"
	"strings"
)

var textMimeTypes = [...]string{
	"application/javascript",
	"application/json",
	"application/x-www-form-urlencoded",
	"text/css",
	"text/html",
	"text/plain",
}

// IsGzipped check if the enconding values refer to a gzipped value
func IsGzipped(encodings ...string) bool {
	for _, encoding := range encodings {
		if strings.ToLower(encoding) == "gzip" {
			return true
		}
	}
	return false
}

// IsText checks if a specified mime type is considered text
func IsText(contentTypes ...string) bool {
	for _, contentType := range contentTypes {
		for _, textMimeType := range textMimeTypes {
			if strings.HasPrefix(strings.ToLower(contentType), textMimeType) {
				return true
			}
		}
	}

	return false
}

// NotFound sends an HTTP 404 response with a text/plain message
func NotFound(req *http.Request, w http.ResponseWriter, path string) {
	if path == "" {
		path = req.URL.String()
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(fmt.Sprintf("Sorry, nothing here for: '%s'\n", path)))
}

// InternalError sends an HTTP 500 with a text/plain message showing the error message
// It also logs the message.
func InternalError(req *http.Request, w http.ResponseWriter, err error) {
	fmt.Printf("Error: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("Sorry, something went wrong: '%s'\n", err)))
}
