package listener

import (
	"sync"

	"github.com/visola/go-proxy/pkg/configuration"
)

// Listener represents a listener that responds to incoming requests
type Listener struct {
	CertificateFile  string               `json:"certificateFile" yaml:"certificateFile"`
	EnabledUpstreams []string             `json:"enabledUpstreams" yaml:"enabledUpstreams"`
	KeyFile          string               `json:"keyFile" yaml:"keyFile"`
	Name             string               `json:"name"`
	Origin           configuration.Origin `json:"-" yaml:"-"`
	Port             int                  `json:"port"`
}

var (
	currentListeners      = make(map[string]*Listener)
	currentListenersMutex sync.Mutex
)

// Listeners return a copy of the listeners in the current state
func Listeners() map[string]Listener {
	result := make(map[string]Listener)
	for k, v := range currentListeners {
		result[k] = *v
	}
	return result
}
