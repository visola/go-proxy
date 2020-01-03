package listener

import (
	"fmt"
	"log"
	"net/http"
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
		listenerToHandle := Listeners()[listernerConfiguration.Port]
		handleRequest(listenerToHandle, req, w)
	})

	if listernerConfiguration.CertificateFile == "" || listernerConfiguration.KeyFile == "" {
		log.Printf("Starting proxy at: http://localhost:%d\n", listernerConfiguration.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", listernerConfiguration.Port), proxyServer); err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Printf("Starting proxy at: https://localhost:%d\n", listernerConfiguration.Port)
	if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", listernerConfiguration.Port), listernerConfiguration.CertificateFile, listernerConfiguration.KeyFile, proxyServer); err != nil {
		log.Fatal(err)
	}
}
