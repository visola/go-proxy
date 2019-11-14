package configuration

import "github.com/visola/go-proxy/pkg/listener"

// ActivateListener activates the lister matched by port and copy the configuration
// of the passed in. Enabled upstreams are loaded from the persisted state.
func ActivateListener(newConfig listener.ListenerConfiguration) {
	currentConfigurationMutex.Lock()
	defer currentConfigurationMutex.Unlock()

	found := false
	for port, l := range currentConfiguration.Listeners {
		if port == newConfig.Port {
			l.Configuration = newConfig
			l.Active = true
			found = true
		}
	}

	if !found {
		currentConfiguration.Listeners[newConfig.Port] = &listener.Listener{
			Active:           true,
			Configuration:    newConfig,
			EnabledUpstreams: make([]string, 0),
		}
	}
}
