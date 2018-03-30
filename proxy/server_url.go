package proxy

import "fmt"

var port = 33443
var isSSL = false

// GetURL returns the full URL to access the proxy server
func GetURL() string {
	protocol := "http"
	if isSSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s:%d", protocol, "localhost", port)
}
