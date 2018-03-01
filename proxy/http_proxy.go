package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/visola/go-proxy/config"
	myhttp "github.com/visola/go-proxy/http"
)

func proxyRequest(req *http.Request, w http.ResponseWriter, mapping config.Mapping) {
	oldPath := req.URL.Path
	newPath := mapping.To + "/" + oldPath[len(mapping.From):]

	newURL, parseErr := url.Parse(newPath)
	if parseErr != nil {
		myhttp.InternalError(req, w, parseErr)
		return
	}

	log.Printf("Proxying '%s' to '%s' for '%s'\n", req.URL.String(), newURL.String(), mapping.Origin)

	newReq, newReqErr := http.NewRequest(req.Method, newURL.String(), req.Body)
	if newReqErr != nil {
		myhttp.InternalError(req, w, newReqErr)
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
		myhttp.InternalError(req, w, respErr)
		return
	}

	defer resp.Body.Close()

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {

			// Fix location headers to point to proxy
			if strings.ToLower(name) == "location" {
				if strings.HasPrefix(value, mapping.To) {
					value = "https://localhost:3443/" + value[len(mapping.To):]
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
			myhttp.InternalError(req, w, readError)
			return
		}

		if bytesRead == 0 {
			break
		}

		w.Write(buffer[:bytesRead])
	}
}
