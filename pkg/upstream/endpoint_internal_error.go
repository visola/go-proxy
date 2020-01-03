package upstream

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/visola/go-proxy/pkg/httputil"
)

func internalServerError(executedURL string, req *http.Request, resp http.ResponseWriter, err error) HandleResult {
	httputil.InternalError(req, resp, err)
	return HandleResult{
		ExecutedURL:  executedURL,
		ResponseBody: ioutil.NopCloser(strings.NewReader(err.Error())),
		ResponseCode: http.StatusInternalServerError,
	}
}
