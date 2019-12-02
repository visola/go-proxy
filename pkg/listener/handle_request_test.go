package listener

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/upstream"
)

func TestHandleRequest(t *testing.T) {
	t.Run("No enabled upstreams", testNoEnabledUpstreams)
	t.Run("Enabled upstream with no endpoints", testWithEnabledUpstreamNoEndpoints)
	t.Run("Enabled upstream with endpoints", testWithMatchingStaticEndpoint)
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

func testWithEnabledUpstreamNoEndpoints(t *testing.T) {
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

func testWithMatchingStaticEndpoint(t *testing.T) {
	dir, err := ioutil.TempDir("", "frontend")
	if err != nil {
		assert.FailNow(t, "Error while creating temp dir", err)
	}

	defer os.RemoveAll(dir)

	fileContent := "Hello world!"
	tempFile := filepath.Join(dir, "hello.txt")
	if err := ioutil.WriteFile(tempFile, []byte(fileContent), 0666); err != nil {
		assert.FailNow(t, "Error while writing to temp file", err)
	}

	upstreamName := "test"
	upstream.AddUpstreams([]upstream.Upstream{
		upstream.Upstream{
			Name: upstreamName,
			StaticEndpoints: []*upstream.StaticEndpoint{
				&upstream.StaticEndpoint{
					BaseEndpoint: upstream.BaseEndpoint{
						From: "/test",
					},
					To: dir,
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

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test/hello.txt", nil)
	resp := httptest.NewRecorder()

	handleRequest(listenerToHandle, req, resp)
	assert.Equal(t, http.StatusOK, resp.Code)

	respBody, respErr := ioutil.ReadAll(resp.Result().Body)
	require.Nil(t, respErr)
	assert.Equal(t, fileContent, string(respBody))
}
