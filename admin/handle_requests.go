package admin

import (
	"log"
	"net/http"

	"github.com/Everbridge/go-proxy/statistics"
)

func handleRequets(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Upgrade") == "websocket" {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println(err)
			return
		}

		statistics.OnRequestProxied(func(proxiedRequest statistics.ProxiedRequest) {
			conn.WriteJSON(proxiedRequest)
		})
	} else {
		responseWithJSON(statistics.GetProxiedRequests(), w, req)
	}
}
