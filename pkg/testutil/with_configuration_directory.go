package testutil

import (
	"os"
	"testing"

	"github.com/visola/go-proxy/pkg/configuration"
)

// WithConfigurationDirectory creates a temporary directory and sets it as the configuration directory
func WithConfigurationDirectory(t *testing.T, test func(*testing.T, string)) func(*testing.T) {
	return WithTempDir(t, func(t *testing.T, tempDir string) {
		previousConfigDir := os.Getenv(configuration.ConfigDirectoryEnvironmentVariable)
		defer func() {
			os.Setenv(configuration.ConfigDirectoryEnvironmentVariable, previousConfigDir)
		}()

		os.Setenv(configuration.ConfigDirectoryEnvironmentVariable, tempDir)

		test(t, tempDir)
	})
}
