package listener

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/configuration"
	"github.com/visola/go-proxy/pkg/testutil"
)

func TestInitialize(t *testing.T) {
	currentListeners = make(map[string]*Listener)
	t.Run("Loads listener correctly", testutil.WithConfigurationDirectory(t, testWithListenerConfigured))

	currentListeners = make(map[string]*Listener)
	t.Run("Ensures at least one listener", testutil.WithConfigurationDirectory(t, testWithNoListenerConfigured))
}

func testWithListenerConfigured(t *testing.T, configDirectory string) {
	// Given a listener is configured
	listenerFile := filepath.Join(configDirectory, listenerSubDirectory, "one.yml")
	l := Listener{
		EnabledUpstreams: make([]string, 0),
		Name:             "One",
		Origin:           configuration.Origin{File: listenerFile},
		Port:             12345,
	}

	saveErr := Save(&l)
	require.Nil(t, saveErr)

	// Clear list of listeners
	currentListeners = make(map[string]*Listener)

	// When initialize gets called
	initErr := Initialize()

	// Then it should not return an error
	require.Nil(t, initErr)

	// And the listener is loaded correctly
	allListeners := Listeners()
	require.Equal(t, 1, len(allListeners))

	loaded := allListeners[l.Name]

	// Avoid flakiness
	l.Origin.LoadedAt = 0
	loaded.Origin.LoadedAt = 0

	assert.Equal(t, l, loaded)
}

func testWithNoListenerConfigured(t *testing.T, configDirectory string) {
	// When initialize gets called
	initErr := Initialize()

	// Then it should not return an error
	require.Nil(t, initErr)

	// And the default listener is configured
	allListeners := Listeners()
	require.Equal(t, 1, len(allListeners))

	l := allListeners[defaultName]
	assert.NotNil(t, l.EnabledUpstreams)
	assert.Equal(t, defaultName, l.Name)
	assert.Equal(t, filepath.Join(configDirectory, listenerSubDirectory, defaultFile), l.Origin.File)
	assert.Equal(t, defaultPort, l.Port)
}
