package listener

import (
	"fmt"
	"log"
	"net/http"
)

func Initialize() {
	configurations := loadConfigurations()
	for _, configuration := range configurations {
		go startServer(configuration)
	}
}

func startServer(configuration ListenerConfiguration) {
	proxyServer := http.NewServeMux()

	proxyServer.HandleFunc("/__listener", func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprintf("%s is up and running!", configuration.Name)))
	})

	if configuration.CertificateFile == "" || configuration.KeyFile == "" {
		log.Printf("Starting proxy at: http://localhost:%d\n", configuration.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", configuration.Port), proxyServer); err != nil {
			log.Fatal(err)
		}
		return
	}

	fmt.Printf("Starting proxy at: https://localhost:%d\n", configuration.Port)
	if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", configuration.Port), configuration.CertificateFile, configuration.KeyFile, proxyServer); err != nil {
		log.Fatal(err)
	}
}
