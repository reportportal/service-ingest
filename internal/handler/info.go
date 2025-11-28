package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func infoRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/info", getInfo)

	return r
}

func getInfo(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"service":"Ingest Service","version":"1.0.0"}`))
}
