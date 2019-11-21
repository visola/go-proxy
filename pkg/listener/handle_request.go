package listener

import (
	"bufio"
	"io"
	"net/http"

	"github.com/visola/go-proxy/pkg/handler"
	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/upstream"
)

func handleRequest(listenerToHandle Listener, req *http.Request, resp http.ResponseWriter) {
	handled := false
	for _, enabledUpstream := range listenerToHandle.EnabledUpstreams {
		candidateUpstream, existsUpstream := upstream.Upstreams()[enabledUpstream]
		if !existsUpstream {
			// This is a weird state but it can happen if mappings changed
			continue
		}

		for _, candidateMapping := range candidateUpstream.Mappings {
			candidateHandler, existsHandler := handler.Handlers[candidateMapping.Type]
			if !existsHandler {
				continue
			}

			if candidateHandler.Matches(candidateMapping, *req) {
				handleResponse := candidateHandler.Handle(candidateMapping, *req)
				handled = true

				for name, values := range handleResponse.Headers {
					for _, value := range values {
						resp.Header().Add(name, value)
					}
				}

				resp.WriteHeader(handleResponse.ResponseCode)

				defer handleResponse.Body.Close()
				if handleResponse.ErrorMessage == "" {
					bufferedReader := bufio.NewReader(handleResponse.Body)
					buffer := make([]byte, 512)
					for {
						bytesRead, readErr := bufferedReader.Read(buffer)
						if readErr == io.EOF {
							break
						}
						resp.Write(buffer[:bytesRead])
					}
				} else {
					resp.Write([]byte(handleResponse.ErrorMessage))
				}
			}
		}
	}

	if !handled {
		httputil.NotFound(req, resp, "Nothing configured to handle this request.")
	}
}
