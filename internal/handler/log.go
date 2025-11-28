package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type logHandler struct{}

func (h logHandler) routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/v2/{projectName}/log", h.createLog)

	return r
}

func (h logHandler) createLog(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}
