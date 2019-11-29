package listener

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/visola/go-proxy/pkg/upstream"
)

func TestHandleRequest(t *testing.T) {
	t.Run("No enabled upstreams", testNoEnabledUpstreams)
	t.Run("Enabled upstream with no mapping", testWithEnabledUpstreamNoMappings)
	t.Run("Enabled upstream with mapping", testWithMatchingMapping)
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
			StaticEndpoints: []upstream.StaticEndpoint{
				upstream.StaticEndpoint{
					BaseEndpoint: upstream.BaseEndpoint{
						From: "/test",
					},
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
	assert.Equal(t, http.StatusOK, resp.Code)
}
