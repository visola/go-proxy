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
)

// StaticEndpoint is an endpoint that responds with file from disk
type StaticEndpoint struct {
	BaseEndpoint
	To string `json:"to"`
}

// Handle handles incoming request matching to files in disk
func (s *StaticEndpoint) Handle(req *http.Request) (int, map[string][]string, io.ReadCloser) {
	if s.Regexp != "" {
		return staticHandleResult(s, replaceRegexp(req.URL.Path, s.To, s.ensureRegexp()), req)
	}

	return staticHandleResult(s, path.Join(s.To, req.URL.Path[len(s.From):]), req)
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

func staticHandleResult(s *StaticEndpoint, pathToReturn string, req *http.Request) (int, map[string][]string, io.ReadCloser) {
	file, err := os.Open(pathToReturn)

	if os.IsNotExist(err) {
		return http.StatusNotFound, nil, ioutil.NopCloser(bytes.NewReader([]byte("File not found: " + pathToReturn)))
	}

	if err != nil {
		return returnError(err)
	}

	headers := map[string][]string{
		"Content-Type": []string{getContentType(file)},
	}

	return http.StatusOK, headers, file
}
