package proxy

import (
	"bytes"
)

type proxyResponse struct {
	executedURL  string
	responseCode int
}

type closeableByteBuffer struct {
	*bytes.Buffer
}

func (bb closeableByteBuffer) Close() error {
	return nil
}
