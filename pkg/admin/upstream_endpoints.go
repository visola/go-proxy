package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/upstream"
)

func registerUpstreamEndpoints(router *mux.Router) {
	router.HandleFunc("/api/upstreams", getUpstreams).Methods(http.MethodGet)
}

func getUpstreams(resp http.ResponseWriter, req *http.Request) {
	upstreamsMap := upstream.Upstreams()
	upstreamsArray := make([]upstream.Upstream, len(upstreamsMap))

	index := 0
	for _, l := range upstreamsMap {
		upstreamsArray[index] = l
		index++
	}

	httputil.RespondWithJSON(upstreamsArray, resp, req)
}
