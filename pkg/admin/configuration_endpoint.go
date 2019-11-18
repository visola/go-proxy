package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/visola/go-proxy/pkg/configuration"
	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/listener"
	"github.com/visola/go-proxy/pkg/upstream"
)

type UpstreamStateChangeResult struct {
	Listener listener.Listener `json:"listener"`
	Upstream upstream.Upstream `json:"upstream"`
}

func registerConfigurationEndpoints(router *mux.Router) {
	router.HandleFunc("/listeners/{listenerPort}/upstreams/{upstreamName}", enableUpstream).Methods(http.MethodPut)
}

func enableUpstream(resp http.ResponseWriter, req *http.Request) {
	listenerPort, portError := strconv.Atoi(mux.Vars(req)["listenerPort"])
	if portError != nil {
		httputil.InternalError(req, resp, portError)
		return
	}

	upstreamName := mux.Vars(req)["upstreamName"]

	upstreamFound, ok := upstream.Upstreams()[upstreamName]
	if !ok {
		httputil.NotFound(req, resp, fmt.Sprintf("Upstrem not found: %s", upstreamName))
		return
	}

	listenerFound, ok := configuration.GetListeners()[listenerPort]
	if !ok {
		httputil.NotFound(req, resp, fmt.Sprintf("Listener not found with port: %d", listenerPort))
		return
	}

	configuration.SetUpstreamState(listenerPort, upstreamName, true)

	result := UpstreamStateChangeResult{
		Listener: listenerFound,
		Upstream: upstreamFound,
	}

	httputil.RespondWithJSON(result, resp, req)
}
