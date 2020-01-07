package listener

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/visola/go-proxy/pkg/event"
	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/upstream"
)

const noHandlerMatchedMessage = "Nothing configured to handle this request."

// Constants to process the request and response bodies
const (
	bufferSize        = 4 * 1024 * 1024        // 4 KBs
	maxStoredBodySize = 5 * 1024 * 1024 * 1024 // 5 MBs
)

func findEndpoints(listenerToHandle Listener) []upstream.Endpoint {
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
	return allEnabledEndpoints
}

func handleRequest(listenerToHandle Listener, req *http.Request, resp http.ResponseWriter) {
	result := newHandleResult(req)
	RequestHandlingChanged(result)

	for _, candidateEndpoint := range findEndpoints(listenerToHandle) {
		if candidateEndpoint.Matches(req) {
			result.Timings.Matched = time.Now().UnixNano()
			RequestHandlingChanged(result)

			statusCode, headers, body := candidateEndpoint.Handle(req)
			result.StatusCode = statusCode
			result.Timings.Handled = time.Now().UnixNano()
			RequestHandlingChanged(result)

			if headers == nil {
				headers = make(map[string][]string)
			}

			result.Response.Headers = headers
			for name, values := range headers {
				for _, value := range values {
					resp.Header().Add(name, value)
				}
			}

			resp.WriteHeader(statusCode)

			bodyBytes, responseError := handleReadCloser(body, resp)
			if responseError == nil {
				result.Response.Body = bodyBytes
			} else {
				result.Error = fmt.Sprintf("Error while reading response: %s", responseError.Error())
			}

			break
		}
	}

	if result.Timings.Matched == 0 {
		httputil.NotFound(req, resp, noHandlerMatchedMessage)
		result.StatusCode = http.StatusNotFound
		result.Response.Body = []byte(noHandlerMatchedMessage)
	}

	result.Timings.Completed = time.Now().UnixNano()
	RequestHandlingChanged(result)
}

func handleReadCloser(readCloser io.ReadCloser, resp http.ResponseWriter) ([]byte, error) {
	defer readCloser.Close()

	responseBytes := make([]byte, 0)
	buffer := make([]byte, bufferSize)
	for {
		bytesRead, readError := readCloser.Read(buffer)

		if readError != nil && readError != io.EOF {
			return nil, readError
		}

		if bytesRead == 0 {
			break
		}

		if len(responseBytes) < maxStoredBodySize {
			responseBytes = append(responseBytes, buffer[:bytesRead]...)
		}
		resp.Write(buffer[:bytesRead])
	}

	return responseBytes, nil
}

func newHandleResult(req *http.Request) event.HandleResult {
	return event.HandleResult{
		ID: uuid.New().String(),
		Request: event.HandleBodyAndHeaders{
			Headers: req.Header,
		},
		Response: event.HandleBodyAndHeaders{},
		Timings: event.HandleTimings{
			Started: time.Now().UnixNano(),
		},
		URL: req.URL.String(),
	}
}
