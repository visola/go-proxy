package admin

import (
	"encoding/json"
	"net/http"

	myhttp "github.com/visola/go-proxy/http"
	"github.com/visola/go-proxy/mapping"
	"github.com/visola/go-proxy/variables"
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

func handleGetPossibleValues(w http.ResponseWriter, req *http.Request) {
	vals, err := variables.GetPossibleValues()
	if err != nil {
		myhttp.InternalError(req, w, err)
		return
	}

	responseWithJSON(vals, w, req)
}

func handleGetSelectedValues(w http.ResponseWriter, req *http.Request) {
	values, valuesError := variables.GetSelectedValues()
	if valuesError != nil {
		myhttp.InternalError(req, w, valuesError)
		return
	}

	responseWithJSON(values, w, req)
}

func handleSelectVariableValue(w http.ResponseWriter, req *http.Request) {
	variable := mux.Vars(req)["variable"]

	decoder := json.NewDecoder(req.Body)
	var value string
	err := decoder.Decode(&value)

	if err != nil {
		myhttp.InternalError(req, w, err)
		return
	}

	setError := variables.SetValue(variable, value)
	if setError != nil {
		myhttp.InternalError(req, w, setError)
		return
	}

	values, valuesError := variables.GetSelectedValues()
	if valuesError != nil {
		myhttp.InternalError(req, w, valuesError)
		return
	}

	responseWithJSON(values, w, req)
}

func handlePutPossibleValues(w http.ResponseWriter, req *http.Request) {
	variable := mux.Vars(req)["variable"]

	decoder := json.NewDecoder(req.Body)
	var passedInValues []string
	err := decoder.Decode(&passedInValues)

	if err != nil {
		myhttp.InternalError(req, w, err)
		return
	}

	vals, err := variables.GetPossibleValues()
	if err != nil {
		myhttp.InternalError(req, w, err)
		return
	}

	vals[variable] = passedInValues
	variables.SetPossibleValues(vals)

	responseWithJSON(vals, w, req)
}

func registerVariablesEndpoints(router *mux.Router) {
	router.HandleFunc("/api/possible-values", handleGetPossibleValues).Methods(http.MethodGet)
	router.HandleFunc("/api/possible-values/{variable}", handlePutPossibleValues).Methods(http.MethodPut)

	router.HandleFunc("/api/values", handleGetSelectedValues).Methods(http.MethodGet)
	router.HandleFunc("/api/values/{variable}", handleSelectVariableValue).Methods(http.MethodPut)

	router.HandleFunc("/api/variables", handleGetVariables).Methods(http.MethodGet)
}
