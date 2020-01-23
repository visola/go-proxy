package listener

import "github.com/visola/go-proxy/pkg/configuration"

// Default values for listener options
const (
	defaultFile = "default.yml"
	defaultName = "Default"
	defaultPort = 33080
)

func ensureListener() error {
	if len(currentListeners) > 0 {
		return nil
	}

	defaultListener := &Listener{
		Name: defaultName,
		Origin: configuration.Origin{},
		Port: defaultPort,
	}

	currentListeners[defaultName] = defaultListener

	return SaveToPersistedState()
}
