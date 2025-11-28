package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type testHandler struct{}

func (h testHandler) routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/v2/{projectName}/item", func(r chi.Router) {
		r.Post("/", h.startRootItem)

		r.Route("/{itemUuid}", func(r chi.Router) {
			r.Post("/", h.startChildItem)
			r.Put("/", h.finishTestItem)
		})
	})

	r.Get("/v1/{projectName}/item/uuid/{itemUuid}", h.getTestItem)

	return r
}

func (h testHandler) startRootItem(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (h testHandler) startChildItem(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (h testHandler) finishTestItem(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (h testHandler) getTestItem(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}
