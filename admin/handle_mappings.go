package admin

import (
	"encoding/json"
	"net/http"

	myhttp "github.com/visola/go-proxy/http"
	"github.com/visola/go-proxy/mapping"
	"github.com/gorilla/mux"
)

func handleGetMappings(w http.ResponseWriter, req *http.Request) {
	respondWithMappings(w, req)
}

func handlePutMapping(w http.ResponseWriter, req *http.Request) {
	mappingID := mux.Vars(req)["mappingID"]

	decoder := json.NewDecoder(req.Body)
	var passedInMapping mapping.Mapping
	err := decoder.Decode(&passedInMapping)

	if err != nil {
		myhttp.InternalError(req, w, err)
		return
	}

	err = mapping.Set(mappingID, passedInMapping)
	if err != nil {
		myhttp.InternalError(req, w, err)
		return
	}

	respondWithMappings(w, req)
}

func handlePutMappings(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var passedInMappings []mapping.Mapping
	err := decoder.Decode(&passedInMappings)

	if err != nil {
		myhttp.InternalError(req, w, err)
		return
	}

	err = mapping.SetAll(passedInMappings)
	if err != nil {
		myhttp.InternalError(req, w, err)
		return
	}

	respondWithMappings(w, req)
}

func registerMappingsEndpoints(router *mux.Router) {
	router.HandleFunc("/api/mappings", handleGetMappings).Methods(http.MethodGet)
	router.HandleFunc("/api/mappings", handlePutMappings).Methods(http.MethodPut)
	router.HandleFunc("/api/mappings/{mappingID}", handlePutMapping).Methods(http.MethodPut)
}

func respondWithMappings(w http.ResponseWriter, req *http.Request) {
	mappings, mappingError := mapping.GetMappings()
	if mappingError != nil {
		myhttp.InternalError(req, w, mappingError)
		return
	}
	responseWithJSON(mappings, w, req)
}
