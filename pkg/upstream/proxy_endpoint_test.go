package upstream

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyEndpointHandle(t *testing.T) {
	responseText := "Hello World!"
	var proxiedRequest *http.Request
	var proxiedBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxiedRequest = r

		defer proxiedRequest.Body.Close()
		var readErr error
		proxiedBody, readErr = ioutil.ReadAll(proxiedRequest.Body)
		require.Nil(t, readErr)

		w.WriteHeader(200)
		w.Write([]byte(responseText))
	}))

	proxyEndpoint := &ProxyEndpoint{
		BaseEndpoint: BaseEndpoint{
			From: "/",
		},
		Headers: map[string][]string{
			"Authorization": []string{"mySecret"},
		},
		To: ts.URL + "/some?param1=one",
	}
	fmt.Println(ts.URL)

	bodyToSend := "Some Body"
	req := httptest.NewRequest(http.MethodPost, "http://nowhere.com/test?param2=two", strings.NewReader(bodyToSend))
	resp := httptest.NewRecorder()

	proxyEndpoint.Handle(req, resp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, responseText, resp.Body.String())

	assert.Equal(t, http.MethodPost, proxiedRequest.Method)
	assert.Equal(t, "/some/test", proxiedRequest.URL.Path)

	require.Equal(t, 1, len(proxiedRequest.Header["Authorization"]))
	assert.Equal(t, "mySecret", proxiedRequest.Header["Authorization"][0])

	assert.Equal(t, "one", proxiedRequest.URL.Query().Get("param1"))
	assert.Equal(t, "two", proxiedRequest.URL.Query().Get("param2"))
	assert.Equal(t, bodyToSend, string(proxiedBody))
}
