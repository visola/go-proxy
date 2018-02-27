package proxy

import (
	"fmt"
	"net/http"
)

func notFound(req *http.Request, w http.ResponseWriter, path string) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Sorry, nothing here for: '%s'\n", path)))
}

func internalError(req *http.Request, w http.ResponseWriter, err error) {
	fmt.Printf("Error: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Sorry, something went wrong: '%s'\n", err)))
}
