package upstream

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUpstreamsFromFile(t *testing.T) {
	fileContent := `
upstreams:
  - name: backend
`

	dir, err := ioutil.TempDir("", "frontend")
	if err != nil {
		assert.FailNow(t, "Error while creating temp dir", err)
	}

	defer os.RemoveAll(dir)

	tempFile := filepath.Join(dir, "go-proxy.yml")
	if err := ioutil.WriteFile(tempFile, []byte(fileContent), 0666); err != nil {
		assert.FailNow(t, "Error while writing to temp file", err)
	}

	loadedUpstreams, err := LoadFromFile(tempFile)
	require.Nil(t, err, "Should load upstreams correctly")

	require.Equal(t, 2, len(loadedUpstreams))

	baseUpstream := loadedUpstreams[0]
	assert.Equal(t, filepath.Base(dir), baseUpstream.Name)
	assert.Equal(t, baseUpstream.Origin.File, tempFile)

	innerUpstream := loadedUpstreams[1]
	assert.Equal(t, "backend", innerUpstream.Name)
	assert.Equal(t, innerUpstream.Origin.File, tempFile)
}

func TestNameFromFilePath(t *testing.T) {
	fileName := nameFromFilePath("/some/crazy/path/backend.yml")
	assert.Equal(t, "backend", fileName)

	parentName := nameFromFilePath("/some/crazy/path/backend/go-proxy.yml")
	assert.Equal(t, "backend", parentName)
}
