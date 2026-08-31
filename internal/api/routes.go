package api

import (
	"net/http"
)

func RegisterHelloRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/hello", HelloHandler)
}
