package upstream

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/testutil"
	"gopkg.in/yaml.v2"
)

func TestCustomDirectory(t *testing.T) {
	testutil.WithConfigurationDirectory(t, func(t *testing.T, tempDir string) {

		customDir1 := "/some/custom/one"
		customDir2 := "/some/custom/two"
		customDir3 := "/some/custom/three"

		AddCustomDirectory(customDir1)
		AddCustomDirectory(customDir2)
		AddCustomDirectory(customDir3)
		assertCustomDirectoriesFile(t, tempDir, customDir1, customDir2, customDir3)

		RemoveCustomDirectory(customDir2)
		assertCustomDirectoriesFile(t, tempDir, customDir1, customDir3)
	})(t)
}

func assertCustomDirectoriesFile(t *testing.T, tempDir string, shouldBe ...string) {
	fileContent, readErr := ioutil.ReadFile(filepath.Join(tempDir, currentCustomDirectoriesFile))
	require.Nil(t, readErr)

	var loadedCustomDirectories []string
	unmarshalError := yaml.Unmarshal(fileContent, &loadedCustomDirectories)
	require.Nil(t, unmarshalError)

	assert.Equal(t, shouldBe, loadedCustomDirectories)
	assert.Equal(t, shouldBe, CustomDirectories())

	loadCustomDirectories()
	assert.Equal(t, shouldBe, CustomDirectories())
}
