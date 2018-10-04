package admin

import (
	"net/http"

	"github.com/Everbridge/go-proxy/configuration"
	myhttp "github.com/Everbridge/go-proxy/http"
	"github.com/gorilla/mux"
)

func handleGetConfigurations(w http.ResponseWriter, req *http.Request) {
	configs, loadErr := configuration.LoadConfiguration()
	if loadErr != nil {
		myhttp.InternalError(req, w, loadErr)
		return
	}

	responseWithJSON(configs, w, req)
}

func registerConfigurationEndpoints(router *mux.Router) {
	router.HandleFunc("/api/configurations", handleGetConfigurations).Methods(http.MethodGet)
}
