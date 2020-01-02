package upstream

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/testutil"
	"gopkg.in/yaml.v2"
)

func TestUmarshalArrayOrString(t *testing.T) {
	content := `
arr:
  - value 1
  - value 2

string: value 3
`

	type testStruct struct {
		Arr    arrayOrString
		String arrayOrString
	}

	testutil.WithTempFile(t, "test.yml", content, func(tempFile *os.File) {
		var yamlContent []byte
		var err error

		yamlContent, err = ioutil.ReadFile(tempFile.Name())
		require.Nil(t, err, "Should read file correctly")

		var yamlStruct testStruct
		err = yaml.Unmarshal(yamlContent, &yamlStruct)
		require.Nil(t, err, "Should unmarshal correctly")

		require.Equal(t, 2, len(yamlStruct.Arr), "Should read two values")
		assert.Equal(t, "value 1", yamlStruct.Arr[0], "Should match correct value")
		assert.Equal(t, "value 2", yamlStruct.Arr[1], "Should match correct value")

		require.Equal(t, 1, len(yamlStruct.String), "Should load one value")
		assert.Equal(t, "value 3", yamlStruct.String[0])
	})
}

func TestFromMapOfArrayOfStrings(t *testing.T) {
	testMap := map[string]arrayOrString{
		"arr": arrayOrString{
			"value 1",
			"value 2",
		},
	}

	convertedMap := FromMapOfArrayOfStrings(testMap)

	values, contains := convertedMap["arr"]
	require.True(t, contains, "Should contain the mapped key")
	require.Equal(t, 2, len(values), "Should contain two value")

	assert.Equal(t, "value 1", values[0], "Should match correct value")
	assert.Equal(t, "value 2", values[1], "Should match correct value")
}
