package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/listener"
	"github.com/visola/go-proxy/pkg/upstream"
)

// UpstreamStateChangeResult struct used to show results of listeners/upstreams endpoint
type UpstreamStateChangeResult struct {
	Listener  listener.Listener   `json:"listener"`
	Upstreams []upstream.Upstream `json:"upstreams"`
}

func registerListenerEndpoints(router *mux.Router) {
	router.HandleFunc("/api/listeners", getListeners).Methods(http.MethodGet)
	router.HandleFunc("/api/listeners/{listenerName}", updateListener).Methods(http.MethodPut)
}

func getListeners(resp http.ResponseWriter, req *http.Request) {
	listenersMap := listener.Listeners()
	listenersArray := make([]listener.Listener, len(listenersMap))

	index := 0
	for _, l := range listenersMap {
		listenersArray[index] = l
		index++
	}

	httputil.RespondWithJSON(listenersArray, resp, req)
}

func updateListener(resp http.ResponseWriter, req *http.Request) {
	listenerName := mux.Vars(req)["listenerName"]

	loadedListener, ok := listener.Listeners()[listenerName]
	if !ok {
		httputil.NotFound(req, resp, fmt.Sprintf("Listener not found: %s", listenerName))
		return
	}

	defer req.Body.Close()
	requestData, readErr := ioutil.ReadAll(req.Body)
	if readErr != nil {
		httputil.InternalError(req, resp, readErr)
	}

	var sentIn listener.Listener
	if err := json.Unmarshal(requestData, &sentIn); err != nil {
		httputil.InternalError(req, resp, err)
		return
	}

	loadedListener.CertificateFile = sentIn.CertificateFile
	loadedListener.EnabledUpstreams = sentIn.EnabledUpstreams
	loadedListener.KeyFile = sentIn.KeyFile
	loadedListener.Port = sentIn.Port

	listener.Save(&loadedListener)
	httputil.RespondWithJSON(loadedListener, resp, req)
}
