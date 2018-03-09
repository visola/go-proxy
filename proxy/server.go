package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/visola/go-proxy/config"
	myhttp "github.com/visola/go-proxy/http"
	"github.com/visola/go-proxy/statistics"
)

var port = 33443

var textMimeTypes = [...]string{
	"application/json",
	"application/x-www-form-urlencoded",
	"text/html",
	"text/plain",
}

// StartProxyServer starts the proxy server
func StartProxyServer() error {
	certFile := os.Getenv("GO_PROXY_CERT_FILE")
	keyFile := os.Getenv("GO_PROXY_CERT_KEY_FILE")

	proxyServer := http.NewServeMux()
	proxyServer.HandleFunc("/", requestHandler)

	fmt.Printf("Starting proxy at: https://localhost:%d\n", port)
	bindError := http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certFile, keyFile, proxyServer)
	return bindError
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	proxiedRequest, configurations, configError := initializeHandling(req, w)

	if configError != nil {
		finalizeWithError(req, w, proxiedRequest, configError)
		return
	}

	config := findConfiguration(req, configurations)

	if config == nil {
		finalizeWithNotFound(req, w, proxiedRequest, "")
		return
	}

	response, handlingError := handleRequest(req, w, config)

	if handlingError == os.ErrNotExist {
		finalizeWithNotFound(req, w, proxiedRequest, response.executedURL)
		return
	}

	if handlingError != nil {
		finalizeWithError(req, w, proxiedRequest, handlingError)
		return
	}

	proxiedRequest.ExecutedURL = response.executedURL
	proxiedRequest.ResponseCode = response.responseCode
	proxiedRequest.EndTime = getTime()
	statistics.AddProxiedRequest(proxiedRequest)
}

func findConfiguration(req *http.Request, configurations []config.Mapping) *config.Mapping {
	for _, config := range configurations {
		if strings.HasPrefix(req.URL.Path, config.From) {
			return &config
		}
	}

	return nil
}

func finalizeWithError(req *http.Request, w http.ResponseWriter, proxiedRequest statistics.ProxiedRequest, err error) {
	myhttp.InternalError(req, w, err)
	proxiedRequest.EndTime = getTime()
	proxiedRequest.Error = err.Error()
	proxiedRequest.ResponseCode = http.StatusNotFound
	statistics.AddProxiedRequest(proxiedRequest)
	return
}

func finalizeWithNotFound(req *http.Request, w http.ResponseWriter, proxiedRequest statistics.ProxiedRequest, executedURL string) {
	myhttp.NotFound(req, w, executedURL)

	proxiedRequest.EndTime = getTime()
	proxiedRequest.ExecutedURL = executedURL
	proxiedRequest.ResponseCode = http.StatusNotFound
	statistics.AddProxiedRequest(proxiedRequest)
}

func getData(req *http.Request) (statistics.HTTPData, error) {
	result := statistics.HTTPData{}
	result.Headers = req.Header

	defer req.Body.Close()
	body, bodyErr := ioutil.ReadAll(req.Body)
	if bodyErr != nil {
		return result, bodyErr
	}

	if len(body) != 0 {
		if isText(req.Header["Content-Type"]) {
			result.Body = string(body)

			// Rewind the body
			req.Body = closeableByteBuffer{bytes.NewBuffer(body)}
		} else {
			result.Body = "Binary"
		}
	}

	return result, nil
}

func getTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func handleRequest(req *http.Request, w http.ResponseWriter, config *config.Mapping) (*proxyResponse, error) {
	if config.Proxy {
		return proxyRequest(req, w, *config)
	}

	return serveStaticFile(req, w, *config)
}

func initializeHandling(req *http.Request, w http.ResponseWriter) (statistics.ProxiedRequest, []config.Mapping, error) {
	reqData, _ := getData(req)

	proxiedRequest := statistics.ProxiedRequest{
		Method:       req.Method,
		RequestedURL: req.URL.String(),
		RequestData:  reqData,
		StartTime:    getTime(),
	}

	config, configError := config.GetConfigurations()
	return proxiedRequest, config, configError
}

func isText(contentTypes []string) bool {
	for _, contentType := range contentTypes {
		for _, textMimeType := range textMimeTypes {
			if strings.HasPrefix(strings.ToLower(contentType), textMimeType) {
				return true
			}
		}
	}

	return false
}
