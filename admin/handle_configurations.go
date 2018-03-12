package admin

import (
	"net/http"

	"github.com/visola/go-proxy/config"
	myhttp "github.com/visola/go-proxy/http"
)

func handleConfigurations(w http.ResponseWriter, req *http.Request) {
	configurations, configError := config.GetConfigurations()
	if configError != nil {
		myhttp.InternalError(req, w, configError)
		return
	}
	responseWithJSON(configurations, w, req)
}
