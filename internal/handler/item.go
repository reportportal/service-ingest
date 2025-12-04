package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type itemHandler struct{}

func (h itemHandler) routesV1() chi.Router {
	r := chi.NewRouter()

	r.Get("/uuid/{itemUuid}", h.getTestItem)

	return r
}

func (h itemHandler) routesV2() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.startRootItem)
	r.Post("/{itemUuid}", h.startChildItem)
	r.Put("/{itemUuid}", h.finishTestItem)

	return r
}

func (h itemHandler) startRootItem(w http.ResponseWriter, r *http.Request) {
	RespondNotImplemented(w, r)
}

func (h itemHandler) startChildItem(w http.ResponseWriter, r *http.Request) {
	RespondNotImplemented(w, r)
}

func (h itemHandler) finishTestItem(w http.ResponseWriter, r *http.Request) {
	RespondNotImplemented(w, r)
}

func (h itemHandler) getTestItem(w http.ResponseWriter, r *http.Request) {
	RespondNotImplemented(w, r)
}
