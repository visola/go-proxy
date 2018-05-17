package admin

import (
	"net/http"

	myhttp "github.com/Everbridge/go-proxy/http"
	"github.com/gorilla/mux"
	"github.com/Everbridge/go-proxy/mapping"
)

func handleGetMappings(w http.ResponseWriter, req *http.Request) {
	respondWithMappings(w, req)
}

func handlePutMapping(w http.ResponseWriter, req *http.Request) {
	mappingID := mux.Vars(req)["mappingID"]
	status := req.URL.Query().Get("active") == "true"

	setStatusErr := mapping.SetStatus(mappingID, status)
	if setStatusErr != nil {
		myhttp.InternalError(req, w, setStatusErr)
		return
	}

	respondWithMappings(w, req)
}

func registerMappingsEndpoints(router *mux.Router) {
	router.HandleFunc("/mappings", handleGetMappings).Methods(http.MethodGet)
	router.HandleFunc("/mappings/{mappingID}", handlePutMapping).Methods(http.MethodPut)
}

func respondWithMappings(w http.ResponseWriter, req *http.Request) {
	mappings, mappingError := mapping.GetMappings()
	if mappingError != nil {
		myhttp.InternalError(req, w, mappingError)
		return
	}
	responseWithJSON(mappings, w, req)
}
