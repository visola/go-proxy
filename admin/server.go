package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/visola/go-proxy/config"

	"github.com/gobuffalo/packr"
)

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

	fmt.Println("Opening admin server at: http://localhost:1234")

	adminServer := http.NewServeMux()
	adminServer.Handle("/", http.FileServer(box))
	adminServer.HandleFunc("/configurations", handleConfigurations)
	return http.ListenAndServe(":1234", adminServer)
}

func handleConfigurations(w http.ResponseWriter, req *http.Request) {
	configurations, configError := config.GetConfigurations()
	if configError != nil {
		// internalError(req, w, configError)
		w.WriteHeader(500)
		return
	}

	json, jsonError := json.Marshal(configurations)
	if jsonError != nil {
		// internalError(req, w, jsonError)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(json))
}
