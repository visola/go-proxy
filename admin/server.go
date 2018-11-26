package admin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Everbridge/go-proxy/configuration"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

	adminServerPort := configuration.GetEnvironment().AdminPort
	fmt.Printf("Opening admin server at: http://localhost:%d\n", adminServerPort)

	adminServer := mux.NewRouter()
	registerConfigurationEndpoints(adminServer)
	registerMappingsEndpoints(adminServer)
	registerVariablesEndpoints(adminServer)
	adminServer.HandleFunc("/api/requests", handleRequets)

	registerReturnToIndexEndpoints(adminServer, box, "/configurations", "/mappings", "/requests", "/variables")

	adminServer.Handle("/{file:.*}", http.FileServer(box))

	return http.ListenAndServe(fmt.Sprintf(":%d", adminServerPort), adminServer)
}

func createReturnIndexHandler(indexHTML string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(indexHTML))
	}
}

func registerReturnToIndexEndpoints(router *mux.Router, box packr.Box, paths ...string) {
	indexHTML := box.String("index.html")
	returnIndexHandler := createReturnIndexHandler(indexHTML)

	for _, path := range paths {
		router.HandleFunc(path, returnIndexHandler)
	}
}
