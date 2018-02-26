package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var configurations []configMapping

func proxyRequest(req *http.Request, w http.ResponseWriter, mapping configMapping) {
	oldPath := req.URL.Path
	newPath := mapping.to + "/" + oldPath[len(mapping.from):]

	newURL, parseErr := url.Parse(newPath)
	if parseErr != nil {
		internalError(req, w, parseErr)
		return
	}

	log.Printf("Proxying '%s' to '%s' for '%s'\n", req.URL.String(), newURL.String(), mapping.origin)

	newReq, newReqErr := http.NewRequest(req.Method, newURL.String(), req.Body)
	if newReqErr != nil {
		internalError(req, w, newReqErr)
		return
	}

	// Copy request headers
	for name, values := range req.Header {
		for _, value := range values {
			newReq.Header.Add(name, value)
		}
	}

	// Copy request cookies
	for _, cookie := range req.Cookies() {
		newReq.AddCookie(cookie)
	}

	client := &http.Client{
		// Do not auto-follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, respErr := client.Do(newReq)

	if respErr != nil {
		internalError(req, w, respErr)
		return
	}

	defer resp.Body.Close()

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {

			// Fix location headers to point to proxy
			if strings.ToLower(name) == "location" {
				if strings.HasPrefix(value, mapping.to) {
					value = "https://localhost:3443/" + value[len(mapping.to):]
				}
			}
			w.Header().Set(name, value)
		}
	}

	// Copy status
	w.WriteHeader(resp.StatusCode)

	buffer := make([]byte, 512)
	for {
		bytesRead, readError := resp.Body.Read(buffer)

		if readError != nil && readError != io.EOF {
			internalError(req, w, readError)
			return
		}

		if bytesRead == 0 {
			break
		}

		w.Write(buffer[:bytesRead])
	}
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	// Find all configurations that match
	path := req.URL.Path
	var served string

	for _, config := range configurations {
		if strings.HasPrefix(path, config.from) {
			if config.proxy {
				proxyRequest(req, w, config)
			} else {
				serveStaticFile(req, w, config)
			}
			return
		}
	}

	if served == "" {
		notFound(req, w, path)
	}
}

func serveStaticFile(req *http.Request, w http.ResponseWriter, mapping configMapping) {
	oldPath := req.URL.Path
	newPath := path.Join(mapping.to, oldPath[len(mapping.from):])

	file, err := os.Open(newPath)

	if err == os.ErrNotExist {
		notFound(req, w, newPath)
		return
	}

	if err != nil {
		internalError(req, w, err)
		return
	}

	defer file.Close()

	log.Printf("Serving '%s' for '%s', from '%s'", newPath, req.URL.Path, mapping.origin)

	buffer := make([]byte, 512)
	loopCount := 0
	for {
		bytesRead, readError := file.Read(buffer)

		if readError != nil && readError != io.EOF {
			internalError(req, w, readError)
		}

		if bytesRead == 0 {
			break
		}

		loopCount++

		contentType := ""
		if loopCount == 1 {
			contentType = mime.TypeByExtension(filepath.Ext(file.Name()))
			if contentType == "" {
				contentType = http.DetectContentType(buffer)
			}
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(buffer[:bytesRead])
	}
}

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

func main() {
	http.HandleFunc("/", requestHandler)

	certFile := os.Getenv("GO_PROXY_CERT_FILE")
	keyFile := os.Getenv("GO_PROXY_CERT_KEY_FILE")

	var configErr error
	configurations, configErr = getConfigurations()
	if configErr != nil {
		panic(configErr)
	}

	for _, config := range configurations {
		fmt.Printf("From: %s\n", config.from)
	}

	startAdminServer()
	err := http.ListenAndServeTLS(":3443", certFile, keyFile, nil)
	if err != nil {
		panic(err)
	}
}
