package listener

import (
	"fmt"
	"log"
	"net/http"
)

// StartListening initializes all the listeners and bind them to the specified ports
func StartListening() {
	for _, l := range currentListeners {
		go startListener(l)
	}
}

func startListener(toStart *Listener) {
	proxyServer := http.NewServeMux()

	proxyServer.HandleFunc("/__listener", func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprintf("%s is up and running!", toStart.Name)))
	})

	proxyServer.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		handleRequest(*currentListeners[toStart.Name], req, w)
	})

	if toStart.CertificateFile == "" || toStart.KeyFile == "" {
		log.Printf("Starting proxy at: http://localhost:%d\n", toStart.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", toStart.Port), proxyServer); err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Printf("Starting proxy at: https://localhost:%d\n", toStart.Port)
	if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", toStart.Port), toStart.CertificateFile, toStart.KeyFile, proxyServer); err != nil {
		log.Fatal(err)
	}
}
