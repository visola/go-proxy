package proxy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/visola/go-proxy/config"
)

var urlToProxy = fmt.Sprintf("https://localhost:%d", port)

func proxyRequest(req *http.Request, w http.ResponseWriter, mapping config.Mapping) (*proxyResponse, error) {
	oldPath := req.URL.Path
	newPath := mapping.To + "/" + oldPath[len(mapping.From):]

	newURL, parseErr := url.Parse(fmt.Sprintf("%s?%s", newPath, req.URL.RawQuery))
	if parseErr != nil {
		return &proxyResponse{
			proxiedTo:    newPath,
			responseCode: http.StatusInternalServerError,
		}, parseErr
	}

	defer req.Body.Close()
	bodyInBytes, readBodyErr := ioutil.ReadAll(req.Body)

	if readBodyErr != nil {
		return &proxyResponse{
			proxiedTo:    newPath,
			responseCode: http.StatusInternalServerError,
		}, readBodyErr
	}

	newReq, newReqErr := http.NewRequest(req.Method, newURL.String(), bytes.NewBuffer(bodyInBytes))
	if newReqErr != nil {
		return &proxyResponse{
			proxiedTo:    newPath,
			responseCode: http.StatusInternalServerError,
		}, newReqErr
	}

	// Copy request headers
	for name, values := range req.Header {
		for _, value := range values {
			newReq.Header.Add(name, value)
		}
	}

	client := &http.Client{
		// Do not auto-follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, respErr := client.Do(newReq)

	if respErr != nil {
		return &proxyResponse{
			proxiedTo:    newPath,
			responseCode: http.StatusInternalServerError,
		}, respErr
	}

	defer resp.Body.Close()

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {

			// Fix location headers to point to proxy
			if strings.ToLower(name) == "location" {
				if strings.HasPrefix(value, mapping.To) {
					value = urlToProxy + value[len(mapping.To):]
				}
			}
			w.Header().Add(name, value)
		}
	}

	// Copy status
	w.WriteHeader(resp.StatusCode)

	buffer := make([]byte, 512)
	for {
		bytesRead, readError := resp.Body.Read(buffer)

		if readError != nil && readError != io.EOF {
			return &proxyResponse{
				proxiedTo:    newPath,
				responseCode: http.StatusInternalServerError,
			}, readError
		}

		if bytesRead == 0 {
			break
		}

		w.Write(buffer[:bytesRead])
	}

	return &proxyResponse{
		proxiedTo:    newPath,
		responseCode: resp.StatusCode,
	}, nil
}
