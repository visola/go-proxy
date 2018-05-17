package main

import (
	"log"

	"github.com/Everbridge/go-proxy/mapping"

	"github.com/Everbridge/go-proxy/admin"
	"github.com/Everbridge/go-proxy/proxy"
)

func main() {
	_, err := mapping.GetMappings()
	if err != nil {
		panic(err)
	}

	go startAdmin()
	startProxy()
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
