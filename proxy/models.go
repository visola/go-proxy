package proxy

import (
	"bytes"
)

type proxyResponse struct {
	body         string
	executedURL  string
	headers      map[string][]string
	responseCode int
}

type closeableByteBuffer struct {
	*bytes.Buffer
}

func (bb closeableByteBuffer) Close() error {
	return nil
}
