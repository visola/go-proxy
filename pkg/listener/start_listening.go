package listener

import (
	"fmt"
	"log"
	"net/http"

	"github.com/visola/go-proxy/pkg/handler"
	"github.com/visola/go-proxy/pkg/upstream"
)

// StartListening initializes all the listeners and bind them to the specified ports
func StartListening(configurations []ListenerConfiguration) {
	for _, configuration := range configurations {
		go startListener(configuration)
	}
}

func startListener(listernerConfiguration ListenerConfiguration) {
	proxyServer := http.NewServeMux()

	proxyServer.HandleFunc("/__listener", func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprintf("%s is up and running!", listernerConfiguration.Name)))
	})

	proxyServer.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		handleRequest(listernerConfiguration, req, w)
	})

	if listernerConfiguration.CertificateFile == "" || listernerConfiguration.KeyFile == "" {
		log.Printf("Starting proxy at: http://localhost:%d\n", listernerConfiguration.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", listernerConfiguration.Port), proxyServer); err != nil {
			log.Fatal(err)
		}
		return
	}

	fmt.Printf("Starting proxy at: https://localhost:%d\n", listernerConfiguration.Port)
	if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", listernerConfiguration.Port), listernerConfiguration.CertificateFile, listernerConfiguration.KeyFile, proxyServer); err != nil {
		log.Fatal(err)
	}
}

func handleRequest(listernerConfiguration ListenerConfiguration, req *http.Request, resp http.ResponseWriter) {
	listenerToHandle := GetListeners()[listernerConfiguration.Port]

	for _, enabledUpstream := range listenerToHandle.EnabledUpstreams {
		candidateUpstream, existsUpstream := upstream.Upstreams()[enabledUpstream]
		if !existsUpstream {
			// This is a weird state but it can happen if mapping files changed
			continue
		}

		for _, candidateMapping := range candidateUpstream.Mappings {
			candidateHandler, existsHandler := handler.Handlers[candidateMapping.Type]
			if !existsHandler {
				continue
			}

			if candidateHandler.Matches(candidateMapping, req) {
				// TODO - Use handler to handle request here
			}
		}
	}
}
