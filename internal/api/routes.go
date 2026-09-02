package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("DELETE /{short_code}", h.Delete)
	mux.HandleFunc("GET /{short_code}", h.Redirect)
}
