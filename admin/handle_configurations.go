package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/visola/go-proxy/config"
	myhttp "github.com/visola/go-proxy/http"
)

func handleGetConfigurations(w http.ResponseWriter, req *http.Request) {
	respondWithConfigurations(w, req)
}

func handlePutConfigurations(w http.ResponseWriter, req *http.Request) {
	mappingID := mux.Vars(req)["mappingID"]
	status := req.URL.Query().Get("active") == "true"

	setStatusErr := config.SetStatus(mappingID, status)
	if setStatusErr != nil {
		myhttp.InternalError(req, w, setStatusErr)
		return
	}

	respondWithConfigurations(w, req)
}

func registerConfigurationEndpoints(router *mux.Router) {
	router.HandleFunc("/configurations", handleGetConfigurations).Methods(http.MethodGet)
	router.HandleFunc("/configurations/{mappingID}", handlePutConfigurations).Methods(http.MethodPut)
}

func respondWithConfigurations(w http.ResponseWriter, req *http.Request) {
	configurations, configError := config.GetConfigurations()
	if configError != nil {
		myhttp.InternalError(req, w, configError)
		return
	}
	responseWithJSON(configurations, w, req)
}
