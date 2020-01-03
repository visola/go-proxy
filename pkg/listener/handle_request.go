package listener

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/upstream"
)

func handleRequest(listenerToHandle Listener, req *http.Request, resp http.ResponseWriter) {
	handled := false
	startTime := time.Now()

	allEnabledEndpoints := make(upstream.Endpoints, 0)
	for _, enabledUpstream := range listenerToHandle.EnabledUpstreams {
		candidateUpstream, existsUpstream := upstream.Upstreams()[enabledUpstream]
		if !existsUpstream {
			// This is a weird state but it can happen if endpoints changed
			continue
		}
		allEnabledEndpoints = append(allEnabledEndpoints, candidateUpstream.Endpoints()...)
	}

	sort.Sort(allEnabledEndpoints)

	for _, candidateEndpoint := range allEnabledEndpoints {
		if candidateEndpoint.Matches(req) {
			handleResult := candidateEndpoint.Handle(req, resp)
			handleResult.HandledBy = candidateEndpoint.Source()
			handleResult.RequestHeaders = req.Header
			handleResult.RequestURL = req.URL.String()
			handleResult.Runtime = time.Now().Sub(startTime).Milliseconds()
			// TODO - Store handleResult somewhere

			fmt.Printf("%s\n", handleResult.HandledBy)
			handled = true
			break
		}
	}

	if !handled {
		httputil.NotFound(req, resp, "Nothing configured to handle this request.")
		// TODO - handle result here
	}
}
