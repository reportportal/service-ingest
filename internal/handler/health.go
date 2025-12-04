package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func healthRouter() chi.Router {
	r := chi.NewRouter()
	r.Head("/", getHealthStatus)
	r.Get("/", getHealthStatus)
	return r
}

func getHealthStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("."))
}
