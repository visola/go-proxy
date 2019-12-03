package listener

import (
	"net/http"

	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/upstream"
)

func handleRequest(listenerToHandle Listener, req *http.Request, resp http.ResponseWriter) {
	handled := false
	for _, enabledUpstream := range listenerToHandle.EnabledUpstreams {
		candidateUpstream, existsUpstream := upstream.Upstreams()[enabledUpstream]
		if !existsUpstream {
			// This is a weird state but it can happen if endpoints changed
			continue
		}

		for _, candidateEndpoint := range candidateUpstream.Endpoints() {
			if candidateEndpoint.Matches(req) {
				candidateEndpoint.Handle(req, resp)
				handled = true
			}
		}
	}

	if !handled {
		httputil.NotFound(req, resp, "Nothing configured to handle this request.")
	}
}
