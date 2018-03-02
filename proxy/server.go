package proxy

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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
	proxiedRequest := statistics.ProxiedRequest{
		StartTime:     time.Now().UnixNano(),
		RequestedPath: req.URL.Path,
	}

	path := req.URL.Path

	configurations, configError := config.GetConfigurations()
	if configError != nil {
		myhttp.InternalError(req, w, configError)
		proxiedRequest.EndTime = time.Now().UnixNano()
		proxiedRequest.Error = configError.Error()
		proxiedRequest.ResponseCode = http.StatusInternalServerError
		statistics.AddProxiedRequest(proxiedRequest)
		return
	}

	var response *proxyResponse
	var handlingError error
	for _, config := range configurations {
		if strings.HasPrefix(path, config.From) {
			proxiedRequest.Mapping = config
			if config.Proxy {
				response, handlingError = proxyRequest(req, w, config)
			} else {
				response, handlingError = serveStaticFile(req, w, config)
			}
			break
		}
	}

	if handlingError == os.ErrNotExist {
		myhttp.NotFound(req, w, response.proxiedTo)

		proxiedRequest.EndTime = time.Now().UnixNano()
		proxiedRequest.Error = handlingError.Error()
		proxiedRequest.ProxiedTo = response.proxiedTo
		proxiedRequest.ResponseCode = http.StatusNotFound
		statistics.AddProxiedRequest(proxiedRequest)
		return
	}

	if handlingError != nil {
		myhttp.InternalError(req, w, configError)

		proxiedRequest.EndTime = time.Now().UnixNano()
		proxiedRequest.Error = handlingError.Error()
		proxiedRequest.ResponseCode = http.StatusInternalServerError
		statistics.AddProxiedRequest(proxiedRequest)
		return
	}

	proxiedRequest.EndTime = time.Now().UnixNano()
	if response == nil {
		proxiedRequest.ResponseCode = http.StatusNotFound
		myhttp.NotFound(req, w, path)
	} else {
		proxiedRequest.ResponseCode = http.StatusOK
	}

	statistics.AddProxiedRequest(proxiedRequest)
}
