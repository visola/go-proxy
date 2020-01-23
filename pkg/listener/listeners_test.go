package listener

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetEnabledUpstreams(t *testing.T) {
	listenerPort := 8000
	listenerName := "test"

	currentListeners = map[string]*Listener{
		listenerName: &Listener{
			Port: listenerPort,
		},
	}

	upstreamNames := []string{"backend", "frontend"}
	SetEnabledUpstreams(listenerName, upstreamNames)
	assertEnabledUpstreamsInListener(t, listenerName, upstreamNames)

	newUpstreamNames := []string{"cluster1", "cluster2"}
	SetEnabledUpstreams(listenerName, newUpstreamNames)
	assertEnabledUpstreamsInListener(t, listenerName, newUpstreamNames)
}

func assertEnabledUpstreamsInListener(t *testing.T, listenerName string, upstreamNames []string) {
	listeners := Listeners()
	require.Equal(t, 1, len(listeners), "Should contain one listener")

	l := listeners[listenerName]
	require.Equal(t, 2, len(l.EnabledUpstreams), "Should have activated upstreams")
	assert.Equal(t, upstreamNames, l.EnabledUpstreams)
}
