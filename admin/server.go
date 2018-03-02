package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/visola/go-proxy/config"
	myhttp "github.com/visola/go-proxy/http"
	"github.com/visola/go-proxy/statistics"

	"github.com/gobuffalo/packr"
)

const (
	adminServerPort = 1234
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

	fmt.Printf("Opening admin server at: http://localhost:%d\n", adminServerPort)

	adminServer := http.NewServeMux()
	adminServer.Handle("/", http.FileServer(box))
	adminServer.HandleFunc("/configurations", handleConfigurations)
	adminServer.HandleFunc("/proxiedRequests", handleProxiedRequests)

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

func handleProxiedRequests(w http.ResponseWriter, req *http.Request) {
	responseWithJSON(statistics.GetProxiedRequests(), w, req)
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
