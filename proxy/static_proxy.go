package proxy

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/visola/go-proxy/config"
)

func serveStaticFile(req *http.Request, w http.ResponseWriter, mapping config.Mapping) (*proxyResponse, error) {
	oldPath := req.URL.Path
	newPath := path.Join(mapping.To, oldPath[len(mapping.From):])

	file, err := os.Open(newPath)

	if err == os.ErrNotExist {
		return &proxyResponse{
			responseCode: http.StatusNotFound,
			proxiedTo:    newPath,
		}, nil
	}

	if err != nil {
		return &proxyResponse{
			responseCode: http.StatusInternalServerError,
			proxiedTo:    newPath,
		}, err
	}

	defer file.Close()

	buffer := make([]byte, 512)
	loopCount := 0
	for {
		bytesRead, readError := file.Read(buffer)

		if readError != nil && readError != io.EOF {
			return &proxyResponse{
				responseCode: http.StatusInternalServerError,
				proxiedTo:    newPath,
			}, readError
		}

		if bytesRead == 0 {
			break
		}

		loopCount++

		contentType := ""
		if loopCount == 1 {
			contentType = mime.TypeByExtension(filepath.Ext(file.Name()))
			if contentType == "" {
				contentType = http.DetectContentType(buffer)
			}
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(buffer[:bytesRead])
	}

	return &proxyResponse{
		responseCode: http.StatusOK,
		proxiedTo:    newPath,
	}, nil
}
