package upstream

import (
	"net/http"

	"github.com/visola/go-proxy/pkg/httputil"
)

func internalServerError(executedURL string, req *http.Request, resp http.ResponseWriter, err error) HandleResult {
	httputil.InternalError(req, resp, err)
	return HandleResult{
		ExecutedURL:  executedURL,
		ResponseBody: []byte(err.Error()),
		ResponseCode: http.StatusInternalServerError,
	}
}
