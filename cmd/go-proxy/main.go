package main

import (
	"log"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"
	"github.com/visola/go-proxy/pkg/admin"
	"github.com/visola/go-proxy/pkg/configuration"
	"github.com/visola/go-proxy/pkg/listener"
	"github.com/visola/go-proxy/pkg/upstream"
)

type CommandLineOptions struct {
	AdminPort int
}

func main() {
	options := parseCommandLineArguments()

	log.Print("Initializing go-proxy...")

	configuration.Initialize()
	listener.Initialize()

	if configDir, err := configuration.GetConfigurationDirectory(); err == nil {
		upstream.Initialize(configDir)
	} else {
		log.Fatalf("Error while finding configuration directory: %v", err)
	}

	startAdminError := admin.StartAdminServer(adminPort(options.AdminPort))
	if startAdminError != nil {
		log.Fatalf("Error while starting admin server: %v", startAdminError)
	}
}

func adminPort(cliPort int) int {
	if cliPort != 0 {
		return cliPort
	}

	adminPort := os.Getenv("GO_PROXY_ADMIN_PORT")
	if adminPort != "" {
		if port, err := strconv.Atoi(adminPort); err == nil {
			return port
		} else {
			log.Fatal("Invalid admin port, not a number: " + adminPort)
		}
	}

	return 3000
}

func parseCommandLineArguments() CommandLineOptions {
	adminPort := flag.IntP("admin-port", "p", 0, "Port to bind the admin server to")

	flag.Parse()

	return CommandLineOptions{
		AdminPort: *adminPort,
	}
}
