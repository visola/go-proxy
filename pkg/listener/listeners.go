package listener

import "sync"

type Listener struct {
	Active           bool                  `json:"active"`
	Configuration    ListenerConfiguration `json:"configuration"`
	EnabledUpstreams []string              `json:"enabledUpstreams"`
}

var (
	currentListeners      = make(map[int]*Listener)
	currentListenersMutex sync.Mutex
)

// ActivateListener activates the lister matched by port and copy the configuration
// passed in. Enabled upstreams are keptd intact if the lister already exist in the
// current map of listeners.
func ActivateListener(newConfig ListenerConfiguration) {
	currentListenersMutex.Lock()
	defer currentListenersMutex.Unlock()

	found := false
	for port, l := range currentListeners {
		if port == newConfig.Port {
			l.Configuration = newConfig
			l.Active = true
			found = true
		}
	}

	if !found {
		currentListeners[newConfig.Port] = &Listener{
			Active:           true,
			Configuration:    newConfig,
			EnabledUpstreams: make([]string, 0),
		}
	}
}

// GetListeners return a copy of the listeners in the current state
func GetListeners() map[int]Listener {
	// TODO - rename GetListeners -> Listeners
	result := make(map[int]Listener)
	for k, v := range currentListeners {
		result[k] = *v
	}
	return result
}

// SetUpstreamState enables or disables an upstream in a specific listener
func SetUpstreamState(listenerPort int, upstreamName string, newState bool) {
	l, exist := currentListeners[listenerPort]

	if !exist {
		return
	}

	currentIndex := -1
	for i, enabledUpstream := range l.EnabledUpstreams {
		if enabledUpstream == upstreamName {
			currentIndex = i
			break
		}
	}

	currentListenersMutex.Lock()
	defer currentListenersMutex.Unlock()

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

// SetListeners reset the current listeners to a new state
func SetListeners(listenersToSet map[int]Listener) {
	currentListenersMutex.Lock()
	defer currentListenersMutex.Unlock()

	currentListeners = make(map[int]*Listener)
	for port, listenerToSet := range listenersToSet {
		currentListeners[port] = &listenerToSet
	}
}

func resetListeners() {
	currentListenersMutex.Lock()
	defer currentListenersMutex.Unlock()

	currentListeners = make(map[int]*Listener)
}
