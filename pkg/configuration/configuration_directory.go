package configuration

import (
	"os"
	"os/user"
)

// ConfigDirectoryEnvironmentVariable stores the name of the environment variable used to load the configuration directory name
const ConfigDirectoryEnvironmentVariable = "GO_PROXY_CONFIGURATION_DIR"

// GetConfigurationDirectory finds where the configuration directory is
func GetConfigurationDirectory() (string, error) {
	fromEnv := os.Getenv(ConfigDirectoryEnvironmentVariable)
	if fromEnv != "" {
		return fromEnv, nil
	}

	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.go-proxy", nil
}
