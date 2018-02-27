package main

import (
	"log"

	"github.com/visola/go-proxy/admin"
	"github.com/visola/go-proxy/proxy"
)

func main() {
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
