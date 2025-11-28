package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type launchHandler struct{}

func (h launchHandler) routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/v1/{projectName}/launch", func(r chi.Router) {
		r.Get("/uuid/{launchUuid}", h.getLaunch)
		r.Put("/{launchId}/update", h.updateLaunch)
	})

	r.Route("/v2/{projectName}/launch", func(r chi.Router) {
		r.Post("/", h.startLaunch)
		r.Post("/merge", h.mergeLaunch)
		r.Put("/{launchUuid}/finish", h.finishLaunch)
	})

	return r
}

func (h launchHandler) startLaunch(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (h launchHandler) finishLaunch(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (h launchHandler) mergeLaunch(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (h launchHandler) getLaunch(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (h launchHandler) updateLaunch(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}
