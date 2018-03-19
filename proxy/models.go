package proxy

import (
	"bytes"
)

type closeableByteBuffer struct {
	*bytes.Buffer
}

func (bb closeableByteBuffer) Close() error {
	return nil
}

type proxyResponse struct {
	body         string
	executedURL  string
	headers      map[string][]string
	responseCode int
}
