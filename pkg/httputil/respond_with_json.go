package httputil

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON returns a 200 OK with a JSON in the body
func RespondWithJSON(data interface{}, w http.ResponseWriter, req *http.Request) {
	json, jsonError := json.Marshal(data)
	if jsonError != nil {
		InternalError(req, w, jsonError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(json))
}
