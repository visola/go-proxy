package upstream

import (
	"io"
	"net/http"
)

const bufferSize = 4 * 1024 * 1024               // 4 KBs
const maxStoredBodySize = 5 * 1024 * 1024 * 1024 // 5 MBs

func handleReadCloser(readCloser io.ReadCloser, executedURL string, req *http.Request, resp http.ResponseWriter) HandleResult {
	defer readCloser.Close()
	responseBytes := make([]byte, 0)
	buffer := make([]byte, bufferSize)

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
