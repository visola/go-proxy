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
static:
  - from: /first
    to: /first/different

  - regexp: /second
    to: /second/different

upstreams:
  - name: backend

    static:
      - from: /first
        to: /first/different

      - regexp: /second
        to: /second/different
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

	loadedUpstreams, err := loadFromFile(tempFile)
	require.Nil(t, err, "Should load upstreams correctly")

	require.Equal(t, 2, len(loadedUpstreams))

	baseUpstream := loadedUpstreams[0]
	assert.Equal(t, filepath.Base(dir), baseUpstream.Name)
	assert.Equal(t, baseUpstream.Origin.File, tempFile)

	assertMappings(t, baseUpstream)

	innerUpstream := loadedUpstreams[1]
	assert.Equal(t, "backend", innerUpstream.Name)
	assert.Equal(t, innerUpstream.Origin.File, tempFile)

	assertMappings(t, innerUpstream)
}

func TestNameFromFilePath(t *testing.T) {
	fileName := nameFromFilePath("/some/crazy/path/backend.yml")
	assert.Equal(t, "backend", fileName)

	parentName := nameFromFilePath("/some/crazy/path/backend/go-proxy.yml")
	assert.Equal(t, "backend", parentName)
}

func assertMappings(t *testing.T, u Upstream) {
	require.Equal(t, 2, len(u.Mappings))

	firstMapping := u.Mappings[0]
	assert.Equal(t, "/first", firstMapping.From)
	assert.Equal(t, MappingTypeStatic, firstMapping.Type)
	assert.Equal(t, u.Name, firstMapping.UpstreamName)

	secondMapping := u.Mappings[1]
	assert.Equal(t, "/second", secondMapping.Regexp)
	assert.Equal(t, MappingTypeStatic, secondMapping.Type)
	assert.Equal(t, u.Name, secondMapping.UpstreamName)
}
