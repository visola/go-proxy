package admin

import (
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
)

func registerStaticEndpoints(router *mux.Router) {
	box := packr.NewBox("../../web/dist") // Relative to this file
	router.Handle("/{file:.*}", http.FileServer(box))
}
