package listener

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/visola/go-proxy/pkg/httputil"
	"github.com/visola/go-proxy/pkg/upstream"
)

const noHandlerMatchedMessage = "Nothing configured to handle this request."

// Constants to process the request and response bodies
const (
	bufferSize        = 4 * 1024 * 1024        // 4 KBs
	maxStoredBodySize = 5 * 1024 * 1024 * 1024 // 5 MBs
)

type handleBodyAndHeaders struct {
	Body    []byte              `json:"body"`
	Headers map[string][]string `json:"headers"`
}

type handleResult struct {
	Error       string `json:"error"`
	ExecutedURL string `json:"executedURL"`
	ID          string `json:"id"`
	HandledBy   string `json:"handledBy"`
	StatusCode  int    `json:"statusCode"`
	URL         string `json:"url"`

	Request  handleBodyAndHeaders `json:"request"`
	Response handleBodyAndHeaders `json:"response"`

	Timings handleTimings `json:"timings"`
}

type handleTimings struct {
	Completed int64 `json:"completed"`
	Handled   int64 `json:"handled"`
	Matched   int64 `json:"matched"`
	Started   int64 `json:"started"`
}

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

	for _, candidateEndpoint := range findEndpoints(listenerToHandle) {
		if candidateEndpoint.Matches(req) {
			result.Timings.Matched = time.Now().UnixNano()

			statusCode, headers, body := candidateEndpoint.Handle(req)
			result.StatusCode = statusCode
			result.Timings.Handled = time.Now().UnixNano()

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

func newHandleResult(req *http.Request) handleResult {
	return handleResult{
		ID: uuid.New().String(),
		Request: handleBodyAndHeaders{
			Headers: req.Header,
		},
		Response: handleBodyAndHeaders{},
		Timings: handleTimings{
			Started: time.Now().UnixNano(),
		},
		URL: req.URL.String(),
	}
}
