package configuration

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigurationDirectory(t *testing.T) {
	// When one requests the configuration directory
	configDir, err := GetConfigurationDirectory()

	// Then it should not return an error
	require.Nil(t, err)

	// And it should be in the user home directory
	homeDir, homeDirErr := os.UserHomeDir()
	require.Nil(t, homeDirErr)
	assert.Equal(t, homeDir, filepath.Dir(configDir))

	// Given the environment variable is set to a directory
	tempDir, err := ioutil.TempDir("", "goproxytest")
	if err != nil {
		assert.FailNow(t, "Error while creating temp dir", err)
	}

	defer os.RemoveAll(tempDir)

	previousConfigDir := os.Getenv(ConfigDirectoryEnvironmentVariable)
	defer func() {
		os.Setenv(ConfigDirectoryEnvironmentVariable, previousConfigDir)
	}()

	os.Setenv(ConfigDirectoryEnvironmentVariable, tempDir)

	// When  one requests the configuration directory
	configDir, err = GetConfigurationDirectory()

	// Then it should not return an error
	require.Nil(t, err)

	// And it should match the environment variable
	assert.Equal(t, tempDir, configDir)
}
