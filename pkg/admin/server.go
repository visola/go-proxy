package admin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// StartAdminServer starts the admin server
func StartAdminServer(port int) error {
	log.Printf("Opening admin server at: http://localhost:%d\n", port)

	adminServer := mux.NewRouter()

	registerListenerEndpoints(adminServer)
	registerRequestsEndpoints(adminServer)
	registerUpstreamEndpoints(adminServer)
	registerWebsocketEndpoint(adminServer)

	registerStaticEndpoints(adminServer)
	adminServer.HandleFunc("/ping", pong).Methods(http.MethodGet)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), adminServer)
}

func pong(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("pong"))
}
