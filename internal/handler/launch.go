package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type launchHandler struct{}

func (rs launchHandler) routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/v1/{projectName}/launch", func(r chi.Router) {
		r.Get("/uuid/{launchUuid}", rs.getLaunch)
		r.Put("/{launchId}/update", rs.updateLaunch)
	})

	r.Route("/v2/{projectName}/launch", func(r chi.Router) {
		r.Post("/", rs.startLaunch)
		r.Post("/merge", rs.mergeLaunch)
		r.Put("/{launchUuid}/finish", rs.finishLaunch)
	})

	return r
}

func (rs launchHandler) startLaunch(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (rs launchHandler) finishLaunch(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (rs launchHandler) mergeLaunch(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (rs launchHandler) getLaunch(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func (rs launchHandler) updateLaunch(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}
