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

// GetListeners return a copy of the listeners in the current configuration
func GetListeners() map[int]listener.Listener {
	result := make(map[int]listener.Listener)
	for k, v := range currentConfiguration.Listeners {
		result[k] = *v
	}
	return result
}

// SetUpstreamState enables or disables an upstream in a specific listener
func SetUpstreamState(listenerPort int, upstreamName string, newState bool) {
	l := currentConfiguration.Listeners[listenerPort]

	currentIndex := -1
	for i, enabledUpstream := range l.EnabledUpstreams {
		if enabledUpstream == upstreamName {
			currentIndex = i
			break
		}
	}

	currentConfigurationMutex.Lock()
	defer currentConfigurationMutex.Unlock()

	if currentIndex == -1 && newState == true {
		// Not in array, add it
		l.EnabledUpstreams = append(l.EnabledUpstreams, upstreamName)
		return
	}

	if currentIndex != -1 && newState == false {
		// In array, remove it
		newValues := l.EnabledUpstreams
		copy(newValues[currentIndex:], newValues[currentIndex+1:])
		l.EnabledUpstreams = newValues[:len(newValues)-1]
	}
}
