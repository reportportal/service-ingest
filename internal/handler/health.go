package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func health() *chi.Mux {
	r := chi.NewRouter()

	r.Head("/health", getHealthStatus)
	r.Get("/health", getHealthStatus)

	return r
}

func getHealthStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("."))
}
