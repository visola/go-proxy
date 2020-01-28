package admin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/upstream"
)

func registerCustomDirectoriesEndpoints(router *mux.Router) {
	router.HandleFunc("/api/custom-directories", addCustomDirectories).Methods(http.MethodPost)
	router.HandleFunc("/api/custom-directories", getCustomDirectories).Methods(http.MethodGet)
	router.HandleFunc("/api/custom-directories", removeCustomDirectory).Methods(http.MethodDelete)
}

func addCustomDirectories(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	requestBodyData, readError := ioutil.ReadAll(req.Body)
	if readError != nil {
		httputil.InternalError(req, resp, readError)
	}

	var customDirectory string
	if err := json.Unmarshal(requestBodyData, &customDirectory); err != nil {
		httputil.InternalError(req, resp, err)
		return
	}

	upstream.AddCustomDirectory(customDirectory)
	httputil.RespondWithJSON(customDirectory, resp, req)
}

func getCustomDirectories(resp http.ResponseWriter, req *http.Request) {
	httputil.RespondWithJSON(upstream.CustomDirectories(), resp, req)
}

func removeCustomDirectory(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	requestBodyData, readError := ioutil.ReadAll(req.Body)
	if readError != nil {
		httputil.InternalError(req, resp, readError)
	}

	var customDirectory string
	if err := json.Unmarshal(requestBodyData, &customDirectory); err != nil {
		httputil.InternalError(req, resp, err)
		return
	}

	upstream.RemoveCustomDirectory(customDirectory)
	httputil.RespondWithJSON(customDirectory, resp, req)
}
