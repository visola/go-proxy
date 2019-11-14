package listener

type Listener struct {
	Active           bool
	Configuration    ListenerConfiguration
	EnabledUpstreams []string
}
