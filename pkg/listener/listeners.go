package listener

import "sync"

// Listener represents a listener that responds to incoming requests
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

// Listeners return a copy of the listeners in the current state
func Listeners() map[int]Listener {
	result := make(map[int]Listener)
	for k, v := range currentListeners {
		result[k] = *v
	}
	return result
}

// SetEnabledUpstreams sets the array of upstreams that are enabled for a specific listener
func SetEnabledUpstreams(listenerPort int, upstreamsToEnable []string) {
	l, exist := currentListeners[listenerPort]

	if !exist {
		return
	}

	currentListenersMutex.Lock()
	defer currentListenersMutex.Unlock()

	l.EnabledUpstreams = make([]string, 0)
	for _, u := range upstreamsToEnable {
		l.EnabledUpstreams = append(l.EnabledUpstreams, u)
	}

	currentListeners[listenerPort] = l
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
