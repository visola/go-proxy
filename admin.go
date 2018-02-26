package main

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/packr"
)

func startAdminServer() {
	box := packr.NewBox("./dist")
	box.Walk(func(fileName string, file packr.File) error {
		fmt.Printf("file: %s\n", fileName)
		return nil
	})

	fmt.Println("Opening admin server at: http://localhost:1234")

	adminServer := http.NewServeMux()
	adminServer.Handle("/", http.FileServer(box))
	http.ListenAndServe(":1234", adminServer)
}
