package listener

type Listener struct {
	Active           bool                  `json:"active"`
	Configuration    ListenerConfiguration `json:"configuration"`
	EnabledUpstreams []string              `json:"enabledUpstreams"`
}
