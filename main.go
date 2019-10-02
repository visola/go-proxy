package main

import (
	"log"
	"os"
	"strconv"

	"github.com/visola/go-proxy/configuration"
	"github.com/visola/go-proxy/mapping"

	"github.com/visola/go-proxy/admin"
	"github.com/visola/go-proxy/proxy"
)

func main() {
	initializeEnvironment()

	_, err := mapping.GetMappings()
	if err != nil {
		panic(err)
	}

	go startAdmin()
	startProxy()
}

func initializeEnvironment() {
	envOptions := make([]configuration.EnvironmentOption, 0)

	certFile := os.Getenv("GO_PROXY_CERT_FILE")
	if certFile != "" {
		envOptions = append(envOptions, configuration.WithCertificateFile(certFile))
	}

	keyFile := os.Getenv("GO_PROXY_CERT_KEY_FILE")
	if keyFile != "" {
		envOptions = append(envOptions, configuration.WithKeyFile(keyFile))
	}

	adminPort := os.Getenv("GO_PROXY_ADMIN_PORT")
	if adminPort != "" {
		if port, err := strconv.Atoi(adminPort); err == nil {
			envOptions = append(envOptions, configuration.WithAdminPort(port))
		} else {
			log.Fatal("Invalid admin port, not a number: " + adminPort)
		}
	}

	proxyPort := os.Getenv("GO_PROXY_PORT")
	if proxyPort != "" {
		if port, err := strconv.Atoi(proxyPort); err == nil {
			envOptions = append(envOptions, configuration.WithProxyPort(port))
		} else {
			log.Fatal("Invalid proxy port, not a number: " + proxyPort)
		}
	}

	configuration.InitializeEnvironment(envOptions...)
}

func startAdmin() {
	log.Println("Starting admin server...")
	adminError := admin.StartAdminServer()
	if adminError != nil {
		panic(adminError)
	}
}

func startProxy() {
	log.Println("Starting proxy server...")
	proxyError := proxy.StartProxyServer()
	if proxyError != nil {
		panic(proxyError)
	}
}
