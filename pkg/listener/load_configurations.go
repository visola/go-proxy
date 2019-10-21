package listener

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// ListenerOptions store the options to create a listener
type ListenerOptions struct {
	CertificateFile string
	KeyFile         string
	Name            string
	Port            int
}

// Default values for listener options
const (
	defaultName = "DEFAULT"
	defaultPort = 33080
)

// Prefixes for listener configuration
const (
	proxyCertificatePrefix = "GO_PROXY_CERT_FILE"
	proxyKeyPrefix         = "GO_PROXY_KEY_FILE"
	proxyPortPrefix        = "GO_PROXY_PORT"
)

func LoadConfigurations() []ListenerOptions {
	listenersByName := make(map[string]*ListenerOptions)

	for _, keyValue := range os.Environ() {
		pair := strings.SplitN(keyValue, "=", 2)
		key := pair[0]
		value := pair[1]

		if strings.HasPrefix(key, proxyPortPrefix) {
			name := getListenerName(key, proxyPortPrefix)
			listener := ensureListener(name, listenersByName)

			if port, err := strconv.Atoi(value); err == nil {
				listener.Port = port
			} else {
				log.Fatalf("Invalid proxy port, not a number: %s = %s", pair[0], pair[1])
			}
		}

		if strings.HasPrefix(key, proxyCertificatePrefix) {
			name := getListenerName(key, proxyCertificatePrefix)
			listener := ensureListener(name, listenersByName)
			listener.CertificateFile = value
		}

		if strings.HasPrefix(key, proxyKeyPrefix) {
			name := getListenerName(key, proxyKeyPrefix)
			listener := ensureListener(name, listenersByName)
			listener.KeyFile = value
		}
	}

	if len(listenersByName) == 0 {
		return []ListenerOptions{
			ListenerOptions{
				Name: defaultName,
				Port: defaultPort,
			},
		}
	}

	listeners := make([]ListenerOptions, len(listenersByName))
	count := 0
	for _, c := range listenersByName {
		listeners[count] = *c
		count++
	}
	return listeners
}

func ensureListener(name string, listenersByName map[string]*ListenerOptions) *ListenerOptions {
	listener := listenersByName[name]
	if listener == nil {
		listener = &ListenerOptions{
			Name: name,
			Port: defaultPort,
		}
		listenersByName[name] = listener
	}
	return listener
}

func getListenerName(key string, prefix string) string {
	if len(key) <= len(prefix) {
		return defaultName
	}

	return key[len(prefix)+1:]
}
