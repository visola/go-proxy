package proxy

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	myhttp "github.com/Everbridge/go-proxy/http"
	"github.com/Everbridge/go-proxy/mapping"
)

func serveStaticFile(req *http.Request, w http.ResponseWriter, match *mapping.MatchResult) (*proxyResponse, error) {
	file, err := os.Open(match.NewPath)

	if err == os.ErrNotExist {
		return &proxyResponse{
			executedURL:  match.NewPath,
			responseCode: http.StatusNotFound,
		}, nil
	}

	if err != nil {
		return &proxyResponse{
			executedURL:  match.NewPath,
			responseCode: http.StatusInternalServerError,
		}, err
	}

	defer file.Close()

	headers := make(map[string][]string)
	buffer := make([]byte, 512)
	contentType := ""
	loopCount := 0
	responseBytes := make([]byte, 0)

	for {
		bytesRead, readError := file.Read(buffer)

		if readError != nil && readError != io.EOF {
			return &proxyResponse{
				executedURL:  match.NewPath,
				responseCode: http.StatusInternalServerError,
			}, readError
		}

		if bytesRead == 0 {
			break
		}

		loopCount++

		if loopCount == 1 {
			contentType = mime.TypeByExtension(filepath.Ext(file.Name()))
			if contentType == "" {
				contentType = http.DetectContentType(buffer)
			}

			headers["Content-Type"] = []string{contentType}
			w.Header().Set("Content-Type", contentType)
		}

		responseBytes = append(responseBytes, buffer[:bytesRead]...)
		w.Write(buffer[:bytesRead])
	}

	bodyString := "Binary Data"
	if myhttp.IsText(contentType) {
		bodyString = string(responseBytes)
	}

	return &proxyResponse{
		body:         bodyString,
		executedURL:  match.NewPath,
		headers:      headers,
		responseCode: http.StatusOK,
	}, nil
}
