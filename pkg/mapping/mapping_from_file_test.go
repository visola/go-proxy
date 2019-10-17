package mapping

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadMappingFromFile(t *testing.T) {
	fileContent := `
proxy:
  - from: /google
    to: https://www.google.com

  - regexp: /test/(.*)
    to: http://localhost:8080/somewhere/$1
    inject:
      headers:
        Authorization: Bearer 1234

static:
  - from: /bing
    to: /some/path

  - regexp: /test_2/(.*)
    to: /another/path/$1
`

	tempFile, err := ioutil.TempFile("", "mappings")
	if err != nil {
		log.Fatal("Error while creating temporary file", err)
	}

	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(fileContent)); err != nil {
		log.Fatal("Error while writing to temp file", err)
	}
	if err := tempFile.Close(); err != nil {
		log.Fatal("Error while closing temporary file", err)
	}

	mappings, err := loadMappingFromFile(tempFile.Name())
	require.Nil(t, err, "Should return error when loading mappings from file")

	assert.Equal(t, 4, len(mappings))

	googleMapping := mappings[0]
	assert.Equal(t, tempFile.Name(), googleMapping.File)
	assert.Equal(t, "/google", googleMapping.From)
	assert.Equal(t, "", googleMapping.Regexp)
	assert.Equal(t, "https://www.google.com", googleMapping.To)
	assert.Equal(t, MappingTypeProxy, googleMapping.Type)

	regexpMapping := mappings[1]
	assert.Equal(t, tempFile.Name(), regexpMapping.File)
	assert.Equal(t, "", regexpMapping.From)
	assert.Equal(t, "/test/(.*)", regexpMapping.Regexp)
	assert.Equal(t, "http://localhost:8080/somewhere/$1", regexpMapping.To)
	assert.Equal(t, MappingTypeProxy, regexpMapping.Type)
	require.Equal(t, 1, len(regexpMapping.Inject.Headers), "Should have copied injected header")
	assert.Equal(t, "Bearer 1234", regexpMapping.Inject.Headers["Authorization"])

	bingStatic := mappings[2]
	assert.Equal(t, tempFile.Name(), bingStatic.File)
	assert.Equal(t, "/bing", bingStatic.From)
	assert.Equal(t, "", bingStatic.Regexp)
	assert.Equal(t, "/some/path", bingStatic.To)
	assert.Equal(t, MappingTypeStatic, bingStatic.Type)

	regexpStatic := mappings[3]
	assert.Equal(t, tempFile.Name(), regexpStatic.File)
	assert.Equal(t, "", regexpStatic.From)
	assert.Equal(t, "/test_2/(.*)", regexpStatic.Regexp)
	assert.Equal(t, "/another/path/$1", regexpStatic.To)
	assert.Equal(t, MappingTypeStatic, regexpStatic.Type)
}
