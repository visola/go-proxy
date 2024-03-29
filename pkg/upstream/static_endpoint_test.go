package upstream

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visola/go-proxy/pkg/testutil"
)

func TestStaticEndpointHandle(t *testing.T) {
	t.Run("404 when file not found", testutil.WithTempDir(t, testFileNotFound))
	t.Run("500 if error", testutil.WithTempDir(t, testReadError))
	t.Run("Matches file using from", testutil.WithTempDir(t, testMatchUsingFrom))
	t.Run("Matches file using regexp", testutil.WithTempDir(t, testMatchUsingRegexp))
}

func testFileNotFound(t *testing.T, tempDir string) {
	staticEndpoint := &StaticEndpoint{
		BaseEndpoint: BaseEndpoint{
			From: "/test",
		},
		To: tempDir,
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test/hello.txt", nil)

	status, executedPath, _, _ := staticEndpoint.Handle(req)
	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, path.Join(tempDir, "/hello.txt"), executedPath)
}

func testMatchUsingFrom(t *testing.T, tempDir string) {
	fileContent := "Hello world!"
	tempFile := filepath.Join(tempDir, "hello.txt")
	if err := ioutil.WriteFile(tempFile, []byte(fileContent), 0666); err != nil {
		assert.FailNow(t, "Error while writing to temp file", err)
	}

	staticEndpoint := &StaticEndpoint{
		BaseEndpoint: BaseEndpoint{
			From: "/test",
		},
		To: tempDir,
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test/hello.txt", nil)

	status, executedPath, headers, body := staticEndpoint.Handle(req)

	assert.Equal(t, http.StatusOK, status)

	respBody, respErr := ioutil.ReadAll(body)
	require.Nil(t, respErr)
	assert.Equal(t, fileContent, string(respBody))

	assert.Equal(t, "text/plain; charset=utf-8", headers["Content-Type"][0])
	assert.Equal(t, tempFile, executedPath)
}

func testMatchUsingRegexp(t *testing.T, tempDir string) {
	fileContent := "Hello world!"
	tempFile := filepath.Join(tempDir, "hello.txt")
	if err := ioutil.WriteFile(tempFile, []byte(fileContent), 0666); err != nil {
		assert.FailNow(t, "Error while writing to temp file", err)
	}

	staticEndpoint := &StaticEndpoint{
		BaseEndpoint: BaseEndpoint{
			Regexp: "/test/(.*).some",
		},
		To: path.Join(tempDir, "$1.txt"),
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test/hello.some", nil)

	status, executedPath, headers, body := staticEndpoint.Handle(req)

	assert.Equal(t, http.StatusOK, status)

	respBody, respErr := ioutil.ReadAll(body)
	require.Nil(t, respErr)
	assert.Equal(t, fileContent, string(respBody))
	assert.Equal(t, tempFile, executedPath)

	assert.Equal(t, "text/plain; charset=utf-8", headers["Content-Type"][0])
}

func testReadError(t *testing.T, tempDir string) {
	fileContent := "Hello world!"
	tempFile := filepath.Join(tempDir, "hello.txt")
	if err := ioutil.WriteFile(tempFile, []byte(fileContent), 0000); err != nil {
		assert.FailNow(t, "Error while writing to temp file", err)
	}

	staticEndpoint := &StaticEndpoint{
		BaseEndpoint: BaseEndpoint{
			From: "/test",
		},
		To: tempDir,
	}

	req := httptest.NewRequest(http.MethodGet, "http://nowhere.com/test/hello.txt", nil)

	status, executedPath, _, _ := staticEndpoint.Handle(req)

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, "", executedPath)
}
