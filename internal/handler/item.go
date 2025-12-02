package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type itemHandler struct{}

func (h itemHandler) routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/v2/{projectName}/item", func(r chi.Router) {
		r.Post("/", h.startRootItem)

		r.Post("/{itemUuid}", h.startChildItem)
		r.Put("/{itemUuid}", h.finishTestItem)
	})

	r.Get("/v1/{projectName}/item/uuid/{itemUuid}", h.getTestItem)

	return r
}

func (h itemHandler) startRootItem(w http.ResponseWriter, r *http.Request) {
	respondNotImplemented(w, r)
}

func (h itemHandler) startChildItem(w http.ResponseWriter, r *http.Request) {
	respondNotImplemented(w, r)
}

func (h itemHandler) finishTestItem(w http.ResponseWriter, r *http.Request) {
	respondNotImplemented(w, r)
}

func (h itemHandler) getTestItem(w http.ResponseWriter, r *http.Request) {
	respondNotImplemented(w, r)
}
