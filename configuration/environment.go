package configuration

import (
	"fmt"
)

// Environment stores information about the runtime
type Environment struct {
	AdminPort       int
	CertificateFile string
	KeyFile         string
	ProxyPort       int
}

// EnvironmentOption represents an option to be applied on the environment
type EnvironmentOption func(*Environment)

// Stores current runtime configuration
var environment *Environment

// GetEnvironment returns the current runtime environment
func GetEnvironment() Environment {
	return *environment
}

// InitializeEnvironment initializes the environment with the specified options
func InitializeEnvironment(options ...EnvironmentOption) {
	if environment != nil {
		panic("Environment is already set.")
	}

	environment = &Environment{
		AdminPort: 1234,
		ProxyPort: 33443,
	}

	fmt.Println("Initializing environment...")
	for _, option := range options {
		option(environment)
	}
}

// WithAdminPort configuration option to setup an Admin server and UI port
func WithAdminPort(port int) EnvironmentOption {
	return func(env *Environment) {
		fmt.Printf("Setting admin port to %d\n", port)
		env.AdminPort = port
	}
}

// WithCertificateFile configuration option to setup the certificate file to use for TLS
func WithCertificateFile(pathToCertFile string) EnvironmentOption {
	return func(env *Environment) {
		fmt.Printf("Setting certificate file to %s\n", pathToCertFile)
		env.CertificateFile = pathToCertFile
	}
}

// WithKeyFile configuration option to setup the key file to use for TLS
func WithKeyFile(pathToKeyFile string) EnvironmentOption {
	return func(env *Environment) {
		fmt.Printf("Setting key file to %s\n", pathToKeyFile)
		env.KeyFile = pathToKeyFile
	}
}

// WithProxyPort configuration option to setup the proxy port
func WithProxyPort(port int) EnvironmentOption {
	return func(env *Environment) {
		fmt.Printf("Setting proxy port to %d\n", port)
		env.ProxyPort = port
	}
}
