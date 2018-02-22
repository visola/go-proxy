package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

var configurations []proxyConfig

func requestHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("server", "go-proxy")

	// Find all configurations that match
	path := req.URL.Path

	for _, config := range configurations {
		for _, staticConfig := range config.Static {
			if strings.HasPrefix(path, staticConfig.From) {
				serveStaticFile(req, w, staticConfig)
				return
			}
		}
	}

	notFound(req, w, path)
}

func serveStaticFile(req *http.Request, w http.ResponseWriter, staticConfig staticConfiguration) {
	oldPath := req.URL.Path
	newPath := path.Join(staticConfig.To, oldPath[len(staticConfig.From):])

	file, err := os.Open(newPath)
	if err != nil {
		notFound(req, w, newPath)
		return
	}

	buffer := make([]byte, 512)
	loopCount := 0
	for {
		bytesRead, readError := file.Read(buffer)

		if readError == io.EOF {
			break
		}

		if readError != nil {
			internalError(req, w, readError)
		}

		loopCount++

		if loopCount == 1 {
			w.Header().Set("Content-Type", http.DetectContentType(buffer))
		}

		w.Write(buffer[:bytesRead])
	}
}

func notFound(req *http.Request, w http.ResponseWriter, path string) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Sorry, nothing here for: '%s'\n", path)))
}

func internalError(req *http.Request, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Sorry, something went wrong: '%s'\n", err)))
}

func main() {
	http.HandleFunc("/", requestHandler)

	certFile := os.Getenv("GO_PROXY_CERT_FILE")
	keyFile := os.Getenv("GO_PROXY_CERT_KEY_FILE")

	var configErr error
	configurations, configErr = getConfigurations()
	if configErr != nil {
		panic(configErr)
	}

	err := http.ListenAndServeTLS(":3443", certFile, keyFile, nil)
	if err != nil {
		panic(err)
	}
}
