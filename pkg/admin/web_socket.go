package admin

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/visola/go-proxy/pkg/event"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type webSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type webSocketRequestHandler struct {
	conn *websocket.Conn
}

// Consume consumes the handle result
func (h *webSocketRequestHandler) Consume(result event.HandleResult) {
	writeErr := h.conn.WriteJSON(webSocketMessage{
		Data: result,
		Type: "request",
	})
	if writeErr != nil {
		h.conn.Close()
		event.RemoveRequestListener(h)
	}
}

func registerWebsocketEndpoint(router *mux.Router) {
	router.HandleFunc("/api/websocket", startWebSocket).Methods(http.MethodGet)
}

func startWebSocket(resp http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	h := &webSocketRequestHandler{
		conn: conn,
	}

	event.AddRequestListener(h)
}
