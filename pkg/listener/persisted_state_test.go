package listener

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/configuration"
	"github.com/visola/go-proxy/pkg/testutil"
)

func TestPersistedState(t *testing.T) {
	testutil.WithConfigurationDirectory(t, func(t *testing.T, tempDir string) {
		listeners := map[string]*Listener{
			"One": &Listener{
				EnabledUpstreams: make([]string, 0),
				Name:             "One",
				Origin: configuration.Origin{
					File: filepath.Join(tempDir, listenerSubDirectory, "/one.yml"),
				},
				Port: 80,
			},
			"Two": &Listener{
				EnabledUpstreams: make([]string, 0),
				Name:             "Two",
				Origin: configuration.Origin{
					File: filepath.Join(tempDir, listenerSubDirectory, "two.yml"),
				},
				Port: 443,
			},
		}

		currentListeners = make(map[string]*Listener)
		for name, val := range listeners {
			currentListeners[name] = val
		}

		saveErr := SaveToPersistedState()
		require.Nil(t, saveErr)
		assertMapListenersEqual(t, listeners, currentListeners)

		loadErr := LoadFromPersistedState()
		require.Nil(t, loadErr)
		assertMapListenersEqual(t, listeners, currentListeners)
	})(t)
}

func assertMapListenersEqual(t *testing.T, expected map[string]*Listener, actual map[string]*Listener) {
	assert.Equal(t, len(expected), len(actual))
	for name, exp := range expected {
		act, exists := actual[name]
		require.True(t, exists)
		exp.Origin.LoadedAt = 0
		act.Origin.LoadedAt = 0
		assert.Equal(t, exp, act)
	}
}
