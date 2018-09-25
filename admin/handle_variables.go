package admin

import (
	"net/http"

	myhttp "github.com/Everbridge/go-proxy/http"
	"github.com/Everbridge/go-proxy/mapping"
	"github.com/Everbridge/go-proxy/variables"
	"github.com/gorilla/mux"
)

func handleGetVariables(w http.ResponseWriter, req *http.Request) {
	allMappings, getError := mapping.GetMappings()
	if getError != nil {
		myhttp.InternalError(req, w, getError)
		return
	}

	varsMap := make(map[string]bool)
	for _, m := range allMappings {
		for _, v := range m.GetVariables() {
			varsMap[v.Name] = true
		}
	}

	vars := make([]string, 0)
	for v := range varsMap {
		vars = append(vars, v)
	}

	responseWithJSON(vars, w, req)
}

func handleGetVariableValues(w http.ResponseWriter, req *http.Request) {
	vals, err := variables.GetPossibleValues()
	if err != nil {
		myhttp.InternalError(req, w, err)
		return
	}

	responseWithJSON(vals, w, req)
}

func registerVariablesEndpoints(router *mux.Router) {
	router.HandleFunc("/api/variables", handleGetVariables).Methods(http.MethodGet)
	router.HandleFunc("/api/variables/values", handleGetVariableValues).Methods(http.MethodGet)
}
