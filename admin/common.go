package admin

import (
	"encoding/json"
	"net/http"

	myhttp "github.com/Everbridge/go-proxy/http"
)

func responseWithJSON(data interface{}, w http.ResponseWriter, req *http.Request) {
	json, jsonError := json.Marshal(data)
	if jsonError != nil {
		myhttp.InternalError(req, w, jsonError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(json))
}
