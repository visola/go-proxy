package context

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// ContextOptions store the options to create a context
type ContextOptions struct {
	CertificateFile string
	KeyFile         string
	Name            string
	Port            int
}

// Default values for context options
const (
	defaultContextName = "DEFAULT"
	defaultPort        = 33080
)

// Prefixes for context configuration
const (
	proxyCertificatePrefix = "GO_PROXY_CERT_FILE"
	proxyKeyPrefix         = "GO_PROXY_KEY_FILE"
	proxyPortPrefix        = "GO_PROXY_PORT"
)

func LoadConfigurations() []ContextOptions {
	contextsByName := make(map[string]*ContextOptions)

	for _, keyValue := range os.Environ() {
		pair := strings.SplitN(keyValue, "=", 2)
		key := pair[0]
		value := pair[1]

		if strings.HasPrefix(key, proxyPortPrefix) {
			name := getContextName(key, proxyPortPrefix)
			context := ensureContext(name, contextsByName)

			if port, err := strconv.Atoi(value); err == nil {
				context.Port = port
			} else {
				log.Fatalf("Invalid proxy port, not a number: %s = %s", pair[0], pair[1])
			}
		}

		if strings.HasPrefix(key, proxyCertificatePrefix) {
			name := getContextName(key, proxyCertificatePrefix)
			context := ensureContext(name, contextsByName)
			context.CertificateFile = value
		}

		if strings.HasPrefix(key, proxyKeyPrefix) {
			name := getContextName(key, proxyKeyPrefix)
			context := ensureContext(name, contextsByName)
			context.KeyFile = value
		}
	}

	if len(contextsByName) == 0 {
		return []ContextOptions{
			ContextOptions{
				Name: defaultContextName,
				Port: defaultPort,
			},
		}
	}

	contexts := make([]ContextOptions, len(contextsByName))
	count := 0
	for _, c := range contextsByName {
		contexts[count] = *c
		count++
	}
	return contexts
}

func ensureContext(name string, contextsByName map[string]*ContextOptions) *ContextOptions {
	context := contextsByName[name]
	if context == nil {
		context = &ContextOptions{
			Name: name,
			Port: defaultPort,
		}
		contextsByName[name] = context
	}
	return context
}

func getContextName(key string, prefix string) string {
	if len(key) <= len(prefix) {
		return defaultContextName
	}

	return key[len(prefix)+1:]
}
