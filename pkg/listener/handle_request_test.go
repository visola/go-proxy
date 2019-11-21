package listener

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/handler"
	"github.com/visola/go-proxy/pkg/upstream"
)

// Dummy handler defaults
const (
	dummyDefaultBody  = "Dummy Response"
	dummyDefaulCode   = http.StatusOK
	dummyDefaultError = ""
	dummyHeaderName   = "Content-type"
	dummyHeaderValue  = "application/json"
)

var (
	dummyBody    string
	dummyCode    int
	dummyError   string
	dummyHeaders map[string][]string
)

type DummyHandler struct{}

func (d DummyHandler) Handle(upstream.Mapping, http.Request) handler.HandleResult {
	return handler.HandleResult{
		Body:         ioutil.NopCloser(strings.NewReader(dummyBody)),
		ErrorMessage: dummyError,
		Headers:      dummyHeaders,
		ResponseCode: dummyCode,
	}
}

func (d DummyHandler) Matches(upstream.Mapping, http.Request) bool {
	return true
}

func TestHandleRequest(t *testing.T) {
	t.Run("No enabled upstreams", testNoEnabledUpstreams)
	t.Run("Enabled upstream with no mapping", testWithEnabledUpstreamNoMappings)
	t.Run("Enabled upstream with mapping", withDummyHandler(testWithMatchingMapping))
	t.Run("Enabled upstream with mapping returns error", withDummyHandler(testWithMatchingMappingReturnsError))
}

func testNoEnabledUpstreams(t *testing.T) {
	listenerToHandle := Listener{
		Configuration: ListenerConfiguration{
			Port: 80,
		},
		EnabledUpstreams: make([]string, 0),
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test", strings.NewReader("Some Body"))
	resp := httptest.NewRecorder()

	handleRequest(listenerToHandle, req, resp)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func testWithEnabledUpstreamNoMappings(t *testing.T) {
	upstreamName := "test"
	upstream.AddUpstreams([]upstream.Upstream{
		upstream.Upstream{
			Name: upstreamName,
		},
	})

	listenerToHandle := Listener{
		Configuration: ListenerConfiguration{
			Port: 80,
		},
		EnabledUpstreams: []string{upstreamName},
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test", strings.NewReader("Some Body"))
	resp := httptest.NewRecorder()

	handleRequest(listenerToHandle, req, resp)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func testWithMatchingMapping(t *testing.T) {
	upstreamName := "test"
	upstream.AddUpstreams([]upstream.Upstream{
		upstream.Upstream{
			Name: upstreamName,
			Mappings: []upstream.Mapping{
				upstream.Mapping{
					Type: "dummy",
				},
			},
		},
	})

	listenerToHandle := Listener{
		Configuration: ListenerConfiguration{
			Port: 80,
		},
		EnabledUpstreams: []string{upstreamName},
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test", strings.NewReader("Some Body"))
	resp := httptest.NewRecorder()

	handleRequest(listenerToHandle, req, resp)
	assert.Equal(t, dummyCode, resp.Code)

	assert.Equal(t, dummyHeaderValue, resp.HeaderMap.Get(dummyHeaderName))

	respBody, respErr := ioutil.ReadAll(resp.Result().Body)
	require.Nil(t, respErr)
	assert.Equal(t, dummyBody, string(respBody))
}

func testWithMatchingMappingReturnsError(t *testing.T) {
	upstreamName := "test"
	upstream.AddUpstreams([]upstream.Upstream{
		upstream.Upstream{
			Name: upstreamName,
			Mappings: []upstream.Mapping{
				upstream.Mapping{
					Type: "dummy",
				},
			},
		},
	})

	listenerToHandle := Listener{
		Configuration: ListenerConfiguration{
			Port: 80,
		},
		EnabledUpstreams: []string{upstreamName},
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test", strings.NewReader("Some Body"))
	resp := httptest.NewRecorder()

	dummyCode = http.StatusInternalServerError
	dummyError = "Something went wrong"

	handleRequest(listenerToHandle, req, resp)

	assert.Equal(t, dummyCode, resp.Code)

	respBody, respErr := ioutil.ReadAll(resp.Result().Body)
	require.Nil(t, respErr)
	assert.Equal(t, dummyError, string(respBody))
}

func withDummyHandler(testFunc func(*testing.T)) func(*testing.T) {
	dummyBody = dummyDefaultBody
	dummyCode = dummyDefaulCode
	dummyError = dummyDefaultError
	dummyHeaders = map[string][]string{
		dummyHeaderName: []string{dummyHeaderValue},
	}

	return func(t *testing.T) {
		oldHandlers := handler.Handlers
		defer func() {
			handler.Handlers = oldHandlers
		}()

		handler.Handlers = map[string]handler.Handler{
			"dummy": DummyHandler{},
		}

		testFunc(t)
	}
}
