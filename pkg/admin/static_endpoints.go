package admin

import (
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

func registerStaticEndpoints(router *mux.Router) {
	box := packr.New("Static Files", "../../web/dist") // Relative to this file
	router.Handle("/{file:.*}", http.FileServer(box))
}
