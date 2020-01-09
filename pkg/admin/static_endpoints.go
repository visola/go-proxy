package admin

import (
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
)

var staticRoots = []string{"/listeners", "/requests"}

func registerStaticEndpoints(router *mux.Router) {
	box := packr.NewBox("../../web/dist") // Relative to this file

	indexHTML := box.String("index.html")
	returnIndexHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(indexHTML))
	}

	for _, path := range staticRoots {
		router.HandleFunc(path, returnIndexHandler)
	}

	router.Handle("/{file:.*}", http.FileServer(box))
}
