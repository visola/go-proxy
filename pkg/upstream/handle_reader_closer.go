package upstream

import (
	"io"
	"net/http"
)

const maxStoredBodySize = 5 * 1024 * 1024 * 1024 // 5 MBs

// TODO - Add tests for this
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

		if len(responseBytes) < maxStoredBodySize {
			responseBytes = append(responseBytes, buffer[:bytesRead]...)
		}
		resp.Write(buffer[:bytesRead])
	}

	return HandleResult{
		ExecutedURL:  executedURL,
		ResponseBody: responseBytes,
	}
}
