package proxy

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	myhttp "github.com/visola/go-proxy/http"
	"github.com/visola/go-proxy/mapping"
)

func proxyRequest(req *http.Request, w http.ResponseWriter, match *mapping.MatchResult) (*proxyResponse, error) {
	mapping := match.Mapping

	newURL, parseErr := url.Parse(fmt.Sprintf("%s?%s", match.NewPath, req.URL.RawQuery))
	if parseErr != nil {
		return &proxyResponse{
			executedURL:  match.NewPath,
			responseCode: http.StatusInternalServerError,
		}, parseErr
	}

	defer req.Body.Close()
	bodyInBytes, readBodyErr := ioutil.ReadAll(req.Body)

	if readBodyErr != nil {
		return &proxyResponse{
			executedURL:  match.NewPath,
			responseCode: http.StatusInternalServerError,
		}, readBodyErr
	}

	newReq, newReqErr := http.NewRequest(req.Method, newURL.String(), bytes.NewBuffer(bodyInBytes))
	if newReqErr != nil {
		return &proxyResponse{
			executedURL:  match.NewPath,
			responseCode: http.StatusInternalServerError,
		}, newReqErr
	}

	// Copy request headers
	for name, values := range req.Header {
		for _, value := range values {
			newReq.Header.Add(name, value)
		}
	}

	// Inject headers if any, overwrite headers passed before
	for name, value := range mapping.Inject.Headers {
		newReq.Header.Set(name, value)
	}

	client := &http.Client{
		// Do not auto-follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, respErr := client.Do(newReq)

	if respErr != nil {
		return &proxyResponse{
			executedURL:  match.NewPath,
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
					value = req.URL.Host + value[len(mapping.To):]
				}
			}
			w.Header().Add(name, value)
		}
	}

	// Copy status
	w.WriteHeader(resp.StatusCode)

	responseBytes := make([]byte, 0)
	buffer := make([]byte, 512)
	for {
		bytesRead, readError := resp.Body.Read(buffer)

		if readError != nil && readError != io.EOF {
			return &proxyResponse{
				executedURL:  match.NewPath,
				responseCode: http.StatusInternalServerError,
			}, readError
		}

		if bytesRead == 0 {
			break
		}

		responseBytes = append(responseBytes, buffer[:bytesRead]...)
		w.Write(buffer[:bytesRead])
	}

	bodyString := "Binary"
	if myhttp.IsText(resp.Header["Content-Type"]...) {
		if myhttp.IsGzipped(resp.Header["Content-Encoding"]...) {
			gzippedReader, _ := gzip.NewReader(bytes.NewReader(responseBytes))
			ungzippedBytes, _ := ioutil.ReadAll(gzippedReader)
			bodyString = string(ungzippedBytes)
		} else {
			bodyString = string(responseBytes)
		}
	}

	return &proxyResponse{
		body:         bodyString,
		executedURL:  match.NewPath,
		headers:      resp.Header,
		responseCode: resp.StatusCode,
	}, nil
}
