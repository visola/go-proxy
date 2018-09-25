package configuration

import "os/user"

// GetConfigurationDirectory returns the directory where all state is stored
func GetConfigurationDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.go-proxy", nil
}
