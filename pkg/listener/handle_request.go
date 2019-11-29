package listener

import (
	"bufio"
	"io"
	"net/http"

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

		for _, candidateEndpoint := range candidateUpstream.Endpoints() {
			if candidateEndpoint.Matches(*req) {
				handleResponse := candidateEndpoint.Handle(*req)
				handled = true

				for name, values := range handleResponse.Headers {
					for _, value := range values {
						resp.Header().Add(name, value)
					}
				}

				if handleResponse.ErrorMessage != "" {
					resp.WriteHeader(http.StatusInternalServerError)
					resp.Write([]byte(handleResponse.ErrorMessage))
				} else {
					resp.WriteHeader(handleResponse.ResponseCode)

					if handleResponse.Body != nil {
						defer handleResponse.Body.Close()

						bufferedReader := bufio.NewReader(handleResponse.Body)
						buffer := make([]byte, 512)
						for {
							bytesRead, readErr := bufferedReader.Read(buffer)
							if readErr == io.EOF {
								break
							}
							resp.Write(buffer[:bytesRead])
						}
					}
				}
			}
		}
	}

	if !handled {
		httputil.NotFound(req, resp, "Nothing configured to handle this request.")
	}
}
