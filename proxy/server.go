package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/visola/go-proxy/config"
	myhttp "github.com/visola/go-proxy/http"
	"github.com/visola/go-proxy/statistics"
)

var port = 33443

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

	match := matchConfiguration(req, configurations)

	if match == nil {
		finalizeWithNotFound(req, w, proxiedRequest, "")
		return
	}

	response, handlingError := handleRequest(req, w, match)

	if handlingError == os.ErrNotExist {
		finalizeWithNotFound(req, w, proxiedRequest, response.executedURL)
		return
	}

	if handlingError != nil {
		finalizeWithError(req, w, proxiedRequest, handlingError)
		return
	}

	proxiedRequest.ResponseData = statistics.HTTPData{
		Body:    response.body,
		Headers: response.headers,
	}
	proxiedRequest.ExecutedURL = response.executedURL
	proxiedRequest.ResponseCode = response.responseCode
	proxiedRequest.EndTime = getTime()
	statistics.AddProxiedRequest(proxiedRequest)
}

func finalizeWithError(req *http.Request, w http.ResponseWriter, proxiedRequest statistics.ProxiedRequest, err error) {
	myhttp.InternalError(req, w, err)
	proxiedRequest.EndTime = getTime()
	proxiedRequest.Error = err.Error()
	proxiedRequest.ResponseCode = http.StatusInternalServerError
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
		if myhttp.IsText(req.Header["Content-Type"]...) {
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

func handleRequest(req *http.Request, w http.ResponseWriter, match *config.MatchResult) (*proxyResponse, error) {
	if match.Mapping.Proxy {
		return proxyRequest(req, w, match)
	}

	return serveStaticFile(req, w, match)
}

func initializeHandling(req *http.Request, w http.ResponseWriter) (statistics.ProxiedRequest, []config.DynamicMapping, error) {
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
