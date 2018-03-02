package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/visola/go-proxy/config"
	myhttp "github.com/visola/go-proxy/http"
	"github.com/visola/go-proxy/statistics"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	adminServerPort = 1234
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// StartAdminServer starts the admin server
func StartAdminServer() error {
	box := packr.NewBox("../dist") // Relative to this file

	fileCount := 0
	box.Walk(func(fileName string, file packr.File) error {
		fileCount++
		return nil
	})

	if fileCount == 0 {
		return errors.New("No files loaded from box")
	}

	fmt.Printf("Opening admin server at: http://localhost:%d\n", adminServerPort)

	adminServer := mux.NewRouter()
	adminServer.HandleFunc("/configurations", handleConfigurations)
	adminServer.HandleFunc("/requests", handleRequets)
	adminServer.Handle("/{file:.*}", http.FileServer(box))

	return http.ListenAndServe(fmt.Sprintf(":%d", adminServerPort), adminServer)
}

func handleConfigurations(w http.ResponseWriter, req *http.Request) {
	configurations, configError := config.GetConfigurations()
	if configError != nil {
		myhttp.InternalError(req, w, configError)
		return
	}
	responseWithJSON(configurations, w, req)
}

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

func responseWithJSON(data interface{}, w http.ResponseWriter, req *http.Request) {
	json, jsonError := json.Marshal(data)
	if jsonError != nil {
		myhttp.InternalError(req, w, jsonError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(json))
}
