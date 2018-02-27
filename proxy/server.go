package proxy

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/visola/go-proxy/config"
)

// StartProxyServer starts the proxy server
func StartProxyServer() error {
	certFile := os.Getenv("GO_PROXY_CERT_FILE")
	keyFile := os.Getenv("GO_PROXY_CERT_KEY_FILE")

	proxyServer := http.NewServeMux()
	proxyServer.HandleFunc("/", requestHandler)

	fmt.Println("Starting proxy at: https://localhost:3443")
	bindError := http.ListenAndServeTLS(":3443", certFile, keyFile, proxyServer)
	return bindError
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	// Find all configurations that match
	path := req.URL.Path
	var served string

	configurations, configError := config.GetConfigurations()
	if configError != nil {
		internalError(req, w, configError)
		return
	}

	for _, config := range configurations {
		if strings.HasPrefix(path, config.From) {
			if config.Proxy {
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
