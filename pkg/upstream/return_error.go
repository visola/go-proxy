package upstream

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func returnError(err error) (int, map[string][]string, io.ReadCloser) {
	return http.StatusInternalServerError, nil, ioutil.NopCloser(bytes.NewReader([]byte(err.Error())))
}
