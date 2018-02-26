package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobuffalo/packr"
)

func handleConfigurations(w http.ResponseWriter, req *http.Request) {
	json, jsonError := json.Marshal(configurations)
	if jsonError != nil {
		internalError(req, w, jsonError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(json))
}

func startAdminServer() {
	box := packr.NewBox("./dist")
	box.Walk(func(fileName string, file packr.File) error {
		fmt.Printf("file: %s\n", fileName)
		return nil
	})

	fmt.Println("Opening admin server at: http://localhost:1234")

	adminServer := http.NewServeMux()
	adminServer.Handle("/", http.FileServer(box))
	adminServer.HandleFunc("/configurations", handleConfigurations)
	http.ListenAndServe(":1234", adminServer)
}
