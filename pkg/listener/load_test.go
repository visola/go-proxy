package listener

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/testutil"
)

func TestLoadFromFile(t *testing.T) {
	testutil.WithConfigurationDirectory(t, func(t *testing.T, tempDir string) {
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
	})(t)
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
