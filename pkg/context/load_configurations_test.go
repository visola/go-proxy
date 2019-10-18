package context

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigurations(t *testing.T) {
	t.Run("Creates default configuration", testCreatesDefaultConfiguration)
	t.Run("Overrides default configuration", testOverridesDefaultConfiguration)
	t.Run("Creates named configuration", testCreatesNamedConfiguration)
}

func testCreatesDefaultConfiguration(t *testing.T) {
	configurations := LoadConfigurations()
	assert.Equal(t, 1, len(configurations))

	defaultConfig := configurations[0]
	assert.Equal(t, defaultContextName, defaultConfig.Name)
	assert.Equal(t, defaultPort, defaultConfig.Port)
	assert.Empty(t, defaultConfig.CertificateFile)
	assert.Empty(t, defaultConfig.KeyFile)
}

func testCreatesNamedConfiguration(t *testing.T) {
	varName := proxyPortPrefix + "_ANOTHER"
	os.Setenv(varName, "10000")
	defer os.Unsetenv(varName)

	configurations := LoadConfigurations()
	assert.Equal(t, 1, len(configurations))

	defaultConfig := configurations[0]
	assert.Equal(t, 10000, defaultConfig.Port)
	assert.Equal(t, "ANOTHER", defaultConfig.Name)
	assert.Empty(t, defaultConfig.CertificateFile)
	assert.Empty(t, defaultConfig.KeyFile)
}

func testOverridesDefaultConfiguration(t *testing.T) {
	os.Setenv(proxyPortPrefix, "10000")
	defer os.Unsetenv(proxyPortPrefix)

	certificateFile := "/some/invalid/path/certificate"
	os.Setenv(proxyCertificatePrefix, certificateFile)
	defer os.Unsetenv(proxyCertificatePrefix)

	keyFile := "/some/invalid/path/key"
	os.Setenv(proxyKeyPrefix, keyFile)
	defer os.Unsetenv(proxyKeyPrefix)

	configurations := LoadConfigurations()
	assert.Equal(t, 1, len(configurations))

	defaultConfig := configurations[0]
	assert.Equal(t, 10000, defaultConfig.Port)
	assert.Equal(t, defaultContextName, defaultConfig.Name)
	assert.Equal(t, certificateFile, defaultConfig.CertificateFile)
	assert.Equal(t, keyFile, defaultConfig.KeyFile)
}
