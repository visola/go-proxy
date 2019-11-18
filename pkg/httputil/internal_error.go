package httputil

import (
	"fmt"
	"net/http"
)

// InternalError sends an HTTP 500 with a text/plain message showing the error message
// It also logs the message.
func InternalError(req *http.Request, w http.ResponseWriter, err error) {
	fmt.Printf("Error: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("Sorry, something went wrong: '%s'\n", err)))
}
