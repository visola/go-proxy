package configuration

import (
	"os"
)

// ConfigDirectoryEnvironmentVariable stores the name of the environment variable used to load the configuration directory name
const ConfigDirectoryEnvironmentVariable = "GO_PROXY_CONFIGURATION_DIR"

// GetConfigurationDirectory finds where the configuration directory is
func GetConfigurationDirectory() (string, error) {
	fromEnv := os.Getenv(ConfigDirectoryEnvironmentVariable)
	if fromEnv != "" {
		return fromEnv, nil
	}

	homeDir, err := os.UserHomeDir()
	return homeDir + "/.go-proxy", err
}
