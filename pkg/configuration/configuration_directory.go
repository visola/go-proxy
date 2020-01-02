package configuration

import (
	"os"
	"os/user"
)

const configDirectoryEnvironmentVariable = "GO_PROXY_CONFIGURATION_DIR"

// GetConfigurationDirectory finds where the configuration directory is
func GetConfigurationDirectory() (string, error) {
	fromEnv := os.Getenv(configDirectoryEnvironmentVariable)
	if fromEnv != "" {
		return fromEnv, nil
	}

	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.go-proxy", nil
}
