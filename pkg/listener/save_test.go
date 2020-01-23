package listener

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/configuration"
)

func TestSave(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "goproxytest")
	if err != nil {
		assert.FailNow(t, "Error while creating temp dir", err)
	}

	defer os.RemoveAll(tempDir)

	previousConfigDir := os.Getenv(configuration.ConfigDirectoryEnvironmentVariable)
	defer func() {
		os.Setenv(configuration.ConfigDirectoryEnvironmentVariable, previousConfigDir)
	}()

	os.Setenv(configuration.ConfigDirectoryEnvironmentVariable, tempDir)

	firstTest := &Listener{
		CertificateFile:  "/path/to/some.crt",
		EnabledUpstreams: []string{"one", "two"},
		KeyFile:          "/path/to/some.key",
		Name:             "My Listener",
		Origin: configuration.Origin{
			File: filepath.Join(tempDir, listenerSubDirectory, "some.yml"),
		},
		Port: 10000,
	}

	testSave(t, firstTest, firstTest.Origin.File)

	secondTest := &Listener{
		CertificateFile:  "/path/to/some.crt",
		EnabledUpstreams: []string{"one", "two"},
		KeyFile:          "/path/to/some.key",
		Name:             "My Listener",
		Port:             10000,
	}

	testSave(t, secondTest, filepath.Join(tempDir, listenerSubDirectory, defaultFile))
}

func testSave(t *testing.T, toSave *Listener, expectedFileName string) {
	saveErr := Save(toSave)
	require.Nil(t, saveErr)

	assert.Equal(t, expectedFileName, toSave.Origin.File)

	loadedContent, loadErr := ioutil.ReadFile(toSave.Origin.File)
	require.Nil(t, loadErr)

	expectedContent := `certificateFile: /path/to/some.crt
enabledUpstreams:
- one
- two
keyFile: /path/to/some.key
name: My Listener
port: 10000
`

	assert.Equal(t, expectedContent, string(loadedContent))
}
