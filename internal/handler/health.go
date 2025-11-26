package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func health() *chi.Mux {
	r := chi.NewRouter()

	r.Head("/health", func(w http.ResponseWriter, r *http.Request) {
		writeHealthResponse(w)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeHealthResponse(w)
	})

	return r
}

func writeHealthResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("."))
}
