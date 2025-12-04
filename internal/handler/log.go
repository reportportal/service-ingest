package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type logHandler struct{}

func (h logHandler) routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.createLog)

	return r
}

func (h logHandler) createLog(w http.ResponseWriter, r *http.Request) {
	respondNotImplemented(w, r)
}
