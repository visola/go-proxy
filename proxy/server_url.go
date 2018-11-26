package proxy

import (
	"fmt"

	"github.com/Everbridge/go-proxy/configuration"
)

var isSSL = false

// GetURL returns the full URL to access the proxy server
func GetURL() string {
	protocol := "http"
	if isSSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s:%d", protocol, "localhost", configuration.GetEnvironment().ProxyPort)
}
