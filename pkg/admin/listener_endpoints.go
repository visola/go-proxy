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

func registerListenerEndpoints(router *mux.Router) {
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

	if _, ok := listener.Listeners()[listenerPort]; !ok {
		httputil.NotFound(req, resp, fmt.Sprintf("Listener not found with port: %d", listenerPort))
		return
	}

	listener.SetUpstreamState(listenerPort, upstreamName, true)
	if err := configuration.SaveToPersistedState(); err != nil {
		httputil.InternalError(req, resp, err)
	}

	result := UpstreamStateChangeResult{
		Listener: listener.Listeners()[listenerPort],
		Upstream: upstreamFound,
	}

	httputil.RespondWithJSON(result, resp, req)
}
