package upstream

import (
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

func getContentType(file *os.File) string {
	contentType := ""

	contentType = mime.TypeByExtension(filepath.Ext(file.Name()))
	if contentType != "" {
		return contentType
	}

	buffer := make([]byte, 512)
	_, readError := file.Read(buffer)
	if readError != nil {
		return ""
	}

	file.Seek(0, 0) // rewind the file
	return http.DetectContentType(buffer)
}

func staticHandleResult(s *StaticEndpoint, pathToReturn string, req *http.Request, resp http.ResponseWriter) HandleResult {
	executedURL := s.UpstreamName + ":" + s.To
	file, err := os.Open(pathToReturn)

	if os.IsNotExist(err) {
		httputil.NotFound(req, resp, "File not found: "+pathToReturn)
		return HandleResult{
			ExecutedURL:  executedURL,
			ResponseBody: ioutil.NopCloser(strings.NewReader("File not found: " + pathToReturn)),
			ResponseCode: http.StatusNotFound,
		}
	}

	if err != nil {
		return internalServerError(executedURL, req, resp, err)
	}

	contentType := getContentType(file)
	resp.Header().Set("Content-Type", contentType)

	handleResult := handleReadCloser(file, executedURL, req, resp)
	handleResult.ResponseCode = http.StatusOK
	handleResult.ResponseHeaders = resp.Header()

	return handleResult
}
