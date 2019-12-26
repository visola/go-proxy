package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/listener"
	"github.com/visola/go-proxy/pkg/upstream"
)

type UpstreamStateChangeResult struct {
	Listener listener.Listener   `json:"listener"`
	Upstream []upstream.Upstream `json:"upstreams"`
}

func registerListenerEndpoints(router *mux.Router) {
	router.HandleFunc("/listeners", getListeners).Methods(http.MethodGet)
	router.HandleFunc("/listeners/{listenerPort}/upstreams", enableUpstream).Methods(http.MethodPut)
}

func enableUpstream(resp http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var upstreamNames []string
	decodeErr := decoder.Decode(&upstreamNames)
	if decodeErr != nil {
		httputil.InternalError(req, resp, decodeErr)
		return
	}

	listenerPort, portError := strconv.Atoi(mux.Vars(req)["listenerPort"])
	if portError != nil {
		httputil.InternalError(req, resp, portError)
		return
	}

	loadedUpstreams := make([]upstream.Upstream, 0)
	for _, upstreamName := range upstreamNames {
		loadedUpstream, ok := upstream.Upstreams()[upstreamName]
		if !ok {
			httputil.NotFound(req, resp, fmt.Sprintf("Upstrem not found: %s", upstreamName))
			return
		}
		loadedUpstreams = append(loadedUpstreams, loadedUpstream)
	}

	if _, ok := listener.Listeners()[listenerPort]; !ok {
		httputil.NotFound(req, resp, fmt.Sprintf("Listener not found with port: %d", listenerPort))
		return
	}

	listener.SetEnabledUpstreams(listenerPort, upstreamNames)

	result := UpstreamStateChangeResult{
		Listener: listener.Listeners()[listenerPort],
		Upstream: loadedUpstreams,
	}

	httputil.RespondWithJSON(result, resp, req)
}

func getListeners(resp http.ResponseWriter, req *http.Request) {
	httputil.RespondWithJSON(listener.Listeners(), resp, req)
}
