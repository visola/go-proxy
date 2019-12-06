package listener

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetUpstreamState(t *testing.T) {
	t.Run("Disable upstream for listener", testDisableUpstreamInListener)
	t.Run("Enable upstream for listener", testEnableUpstreamInListener)
}

func testDisableUpstreamInListener(t *testing.T) {
	upstreamName := "backend"
	listenerPort := 8000

	testEnableUpstreamInListener(t)

	SetUpstreamState(listenerPort, upstreamName, false)

	listeners := Listeners()

	require.Equal(t, 1, len(listeners), "Should contain one listener")

	l := listeners[listenerPort]
	require.Empty(t, l.EnabledUpstreams, "Should have no active upstream")
}

func testEnableUpstreamInListener(t *testing.T) {
	resetListeners()

	upstreamName := "backend"
	listenerPort := 8000
	ActivateListener(ListenerConfiguration{
		Port: listenerPort,
	})

	SetUpstreamState(listenerPort, upstreamName, true)

	listeners := Listeners()

	require.Equal(t, 1, len(listeners), "Should contain one listener")

	l := listeners[listenerPort]
	require.Equal(t, 1, len(l.EnabledUpstreams), "Should have one active upstream")
	assert.Equal(t, upstreamName, l.EnabledUpstreams[0])
}
