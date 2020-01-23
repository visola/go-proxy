package listener

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/testutil"
	"github.com/visola/go-proxy/pkg/upstream"
)

func TestHandleRequest(t *testing.T) {
	t.Run("No enabled upstreams", testNoEnabledUpstreams)
	t.Run("Enabled upstream with no endpoints", testWithEnabledUpstreamNoEndpoints)
	t.Run("Enabled upstream with endpoints", testutil.WithTempDir(t, testWithMatchingStaticEndpoint))
}

func testNoEnabledUpstreams(t *testing.T) {
	listenerToHandle := Listener{
		EnabledUpstreams: make([]string, 0),
		Port:             80,
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
		EnabledUpstreams: []string{upstreamName},
		Port:             80,
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test", strings.NewReader("Some Body"))
	resp := httptest.NewRecorder()

	handleRequest(listenerToHandle, req, resp)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func testWithMatchingStaticEndpoint(t *testing.T, tempDir string) {
	fileContent := "Hello world!"
	tempFile := filepath.Join(tempDir, "hello.txt")
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
					To: tempDir,
				},
			},
		},
	})

	listenerToHandle := Listener{
		EnabledUpstreams: []string{upstreamName},
		Port:             80,
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test/hello.txt", nil)
	resp := httptest.NewRecorder()

	handleRequest(listenerToHandle, req, resp)
	assert.Equal(t, http.StatusOK, resp.Code)

	respBody, respErr := ioutil.ReadAll(resp.Result().Body)
	require.Nil(t, respErr)
	assert.Equal(t, fileContent, string(respBody))
}
