package upstream

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func handleReadCloser(readCloser io.ReadCloser, executedURL string, req *http.Request, resp http.ResponseWriter) HandleResult {
	defer readCloser.Close()
	responseBytes := make([]byte, 0)
	buffer := make([]byte, 512)

	for {
		bytesRead, readError := readCloser.Read(buffer)

		if readError != nil && readError != io.EOF {
			return internalServerError(executedURL, req, resp, readError)
		}

		if bytesRead == 0 {
			break
		}

		responseBytes = append(responseBytes, buffer[:bytesRead]...)
		resp.Write(buffer[:bytesRead])
	}

	return HandleResult{
		ExecutedURL:  executedURL,
		ResponseBody: ioutil.NopCloser(bytes.NewReader(responseBytes)),
	}
}
