package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/listener"
)

func registerRequestsEndpoints(router *mux.Router) {
	router.HandleFunc("/requests", getRequests).Methods(http.MethodGet)
}

func getRequests(resp http.ResponseWriter, req *http.Request) {
	httputil.RespondWithJSON(listener.GetRequests(), resp, req)
}
