package listener

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
	configurations := loadConfigurations()
	assert.Equal(t, 1, len(configurations))

	defaultConfig := configurations[0]
	assert.Equal(t, defaultName, defaultConfig.Name)
	assert.Equal(t, defaultPort, defaultConfig.Port)
	assert.Empty(t, defaultConfig.CertificateFile)
	assert.Empty(t, defaultConfig.KeyFile)
}

func testCreatesNamedConfiguration(t *testing.T) {
	portVar := proxyPortPrefix + "_ANOTHER"
	os.Setenv(portVar, "10000")
	defer os.Unsetenv(portVar)

	certificateFile := "/another/invalid/path/certificate"
	certFileVar := proxyCertificatePrefix + "_ANOTHER"
	os.Setenv(certFileVar, certificateFile)
	defer os.Unsetenv(certFileVar)

	keyFile := "/another/invalid/path/key"
	keyFileVar := proxyKeyPrefix + "_ANOTHER"
	os.Setenv(keyFileVar, keyFile)
	defer os.Unsetenv(keyFileVar)

	configurations := loadConfigurations()
	assert.Equal(t, 1, len(configurations))

	namedConfig := configurations[0]
	assert.Equal(t, 10000, namedConfig.Port)
	assert.Equal(t, "ANOTHER", namedConfig.Name)
	assert.Equal(t, certificateFile, namedConfig.CertificateFile)
	assert.Equal(t, keyFile, namedConfig.KeyFile)
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

	configurations := loadConfigurations()
	assert.Equal(t, 1, len(configurations))

	defaultConfig := configurations[0]
	assert.Equal(t, 10000, defaultConfig.Port)
	assert.Equal(t, defaultName, defaultConfig.Name)
	assert.Equal(t, certificateFile, defaultConfig.CertificateFile)
	assert.Equal(t, keyFile, defaultConfig.KeyFile)
}
