package listener

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetEnabledUpstreams(t *testing.T) {
	listenerPort := 8000

	ActivateListener(ListenerConfiguration{
		Port: listenerPort,
	})

	upstreamNames := []string{"backend", "frontend"}
	SetEnabledUpstreams(listenerPort, upstreamNames)
	assertEnabledUpstreamsInListener(t, listenerPort, upstreamNames)

	newUpstreamNames := []string{"cluster1", "cluster2"}
	SetEnabledUpstreams(listenerPort, newUpstreamNames)
	assertEnabledUpstreamsInListener(t, listenerPort, newUpstreamNames)
}

func assertEnabledUpstreamsInListener(t *testing.T, listenerPort int, upstreamNames []string) {
	listeners := Listeners()
	require.Equal(t, 1, len(listeners), "Should contain one listener")

	l := listeners[listenerPort]
	require.Equal(t, 2, len(l.EnabledUpstreams), "Should have activated upstreams")
	assert.Equal(t, upstreamNames, l.EnabledUpstreams)
}
