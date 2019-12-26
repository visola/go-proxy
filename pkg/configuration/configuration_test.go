package configuration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/listener"
	"github.com/visola/go-proxy/pkg/testutil"
	"github.com/visola/go-proxy/pkg/upstream"
)

func TestPersistedState(t *testing.T) {
	t.Run("Saves and loads correctly", testutil.WithTempDir(t, testSavesAndLoadsCorrecty))
}

func testSavesAndLoadsCorrecty(t *testing.T, tempDir string) {
	previousConfigDir := os.Getenv(configDirectoryEnvironmentVariable)
	defer func() {
		os.Setenv(configDirectoryEnvironmentVariable, previousConfigDir)
	}()

	os.Setenv(configDirectoryEnvironmentVariable, tempDir)

	// Given an upstream
	testUpstream := upstream.Upstream{
		Name: "Test Upstream",
	}
	upstream.AddUpstreams([]upstream.Upstream{testUpstream})

	// And a listener with an active upstream
	testListenerConfig := listener.ListenerConfiguration{
		Name: "Test Listener",
		Port: 10000,
	}
	listener.ActivateListener(testListenerConfig)
	listener.SetEnabledUpstreams(testListenerConfig.Port, []string{testUpstream.Name})

	// When I save to persisted state
	saveErr := SaveToPersistedState()
	require.Nil(t, saveErr, "Should save correctly")

	// Then loading the persisted state should work
	loadErr := LoadFromPersistedState()
	require.Nil(t, loadErr, "Should load correctly")

	// And should contain the saved listener
	allListeners := listener.Listeners()
	assert.Equal(t, 1, len(allListeners))

	// And listener should have correct information
	loadedTestListener, listenerExists := allListeners[testListenerConfig.Port]
	require.True(t, listenerExists, "Should load correct listeners")
	assert.Equal(t, testListenerConfig.Name, loadedTestListener.Configuration.Name)
	assert.Equal(t, testListenerConfig.Port, loadedTestListener.Configuration.Port)

	// And listener should have correctly enabled upstream
	require.Equal(t, 1, len(loadedTestListener.EnabledUpstreams), "Should have one active upstream")
	assert.Equal(t, testUpstream.Name, loadedTestListener.EnabledUpstreams[0])
}
