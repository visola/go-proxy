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

func TestLoadFromFile(t *testing.T) {
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

	contentWithName := `
certificateFile: /path/to/some.crt
enabledUpstreams:
  - one
  - two
keyFile: /path/to/some.key
name: My Listener
port: 10000
`

	testFile(t, filepath.Join(tempDir, "another.yml"), contentWithName, "My Listener")

	contentWithoutName := `
certificateFile: /path/to/some.crt
enabledUpstreams:
  - one
  - two
keyFile: /path/to/some.key
port: 10000
`

	testFile(t, filepath.Join(tempDir, "SomeName.yml"), contentWithoutName, "SomeName")
}

func testFile(t *testing.T, tempFile string, content string, expectedName string) {
	ioutil.WriteFile(tempFile, []byte(content), 0666)

	loadedListener, loadErr := loadFromFile(tempFile)

	require.Nil(t, loadErr)

	assert.Equal(t, "/path/to/some.crt", loadedListener.CertificateFile)
	assert.Equal(t, "/path/to/some.key", loadedListener.KeyFile)
	assert.Equal(t, expectedName, loadedListener.Name)
	assert.Equal(t, 10000, loadedListener.Port)

	require.Equal(t, 2, len(loadedListener.EnabledUpstreams))
	assert.Equal(t, "one", loadedListener.EnabledUpstreams[0])
	assert.Equal(t, "two", loadedListener.EnabledUpstreams[1])

	assert.Equal(t, tempFile, loadedListener.Origin.File)
}
