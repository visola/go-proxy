package upstream

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/visola/go-proxy/pkg/httputil"
)

// StaticEndpoint is an endpoint that responds with file from disk
type StaticEndpoint struct {
	BaseEndpoint
	To string `json:"to"`
}

// Handle handles incoming request matching to files in disk
func (s *StaticEndpoint) Handle(req *http.Request, resp http.ResponseWriter) HandleResult {
	if s.Regexp != "" {
		return staticHandleResult(s, replaceRegexp(req.URL.Path, s.To, s.ensureRegexp()), req, resp)
	}

	return staticHandleResult(s, path.Join(s.To, req.URL.Path[len(s.From):]), req, resp)
}

func staticHandleResult(s *StaticEndpoint, pathToReturn string, req *http.Request, resp http.ResponseWriter) HandleResult {
	executedURL := s.UpstreamName + ":" + s.To
	file, err := os.Open(pathToReturn)

	if os.IsNotExist(err) {
		httputil.NotFound(req, resp, "File not found: "+pathToReturn)
		return HandleResult{
			Body:         ioutil.NopCloser(strings.NewReader("File not found: " + pathToReturn)),
			ExecutedURL:  executedURL,
			ResponseCode: http.StatusNotFound,
		}
	}

	if err != nil {
		httputil.InternalError(req, resp, err)
		return HandleResult{
			Body:         ioutil.NopCloser(strings.NewReader(err.Error())),
			ExecutedURL:  executedURL,
			ResponseCode: http.StatusInternalServerError,
		}
	}

	defer file.Close()

	headers := make(map[string][]string)
	buffer := make([]byte, 512)
	contentType := ""
	responseBytes := make([]byte, 0)

	for {
		bytesRead, readError := file.Read(buffer)

		if readError != nil && readError != io.EOF {
			return HandleResult{
				Body:         ioutil.NopCloser(strings.NewReader(err.Error())),
				ExecutedURL:  executedURL,
				ResponseCode: http.StatusInternalServerError,
			}
		}

		if bytesRead == 0 {
			break
		}

		if contentType == "" {
			contentType = mime.TypeByExtension(filepath.Ext(file.Name()))
			if contentType == "" {
				contentType = http.DetectContentType(buffer)
			}

			headers["Content-Type"] = []string{contentType}
			resp.Header().Set("Content-Type", contentType)
		}

		responseBytes = append(responseBytes, buffer[:bytesRead]...)
		resp.Write(buffer[:bytesRead])
	}

	return HandleResult{
		Body:         ioutil.NopCloser(bytes.NewReader(responseBytes)),
		ExecutedURL:  executedURL,
		Headers:      headers,
		ResponseCode: http.StatusOK,
	}
}
