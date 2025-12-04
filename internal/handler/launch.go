package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/reportportal/service-ingest/internal/service"
)

type launchHandler struct {
	service *service.LaunchService
}

func (h launchHandler) routesV1() chi.Router {
	r := chi.NewRouter()

	r.Get("/uuid/{launchUuid}", h.getLaunch)
	r.Put("/{launchId}/update", h.updateLaunch)

	return r
}

func (h launchHandler) routesV2() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.startLaunch)
	r.Post("/merge", h.mergeLaunch)
	r.Put("/{launchUuid}/finish", h.finishLaunch)

	return r
}

func (h launchHandler) startLaunch(w http.ResponseWriter, r *http.Request) {
	var rq StartLaunchRQ
	if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if err := validate.Struct(rq); err != nil {
		respondValidationError(w, err)
		return
	}

	projectName := chi.URLParam(r, "projectName")
	resp, err := h.service.StartLaunch(rq.toLaunchModel(), projectName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to start launch")
		return
	}

	respondJSON(w, http.StatusCreated, StartLaunchRS{
		UUID:   resp.UUID,
		Number: resp.Number,
	})
}

func (h launchHandler) finishLaunch(w http.ResponseWriter, r *http.Request) {
	respondNotImplemented(w, r)
}

func (h launchHandler) mergeLaunch(w http.ResponseWriter, r *http.Request) {
	respondNotImplemented(w, r)
}

func (h launchHandler) getLaunch(w http.ResponseWriter, r *http.Request) {
	respondNotImplemented(w, r)
}

func (h launchHandler) updateLaunch(w http.ResponseWriter, r *http.Request) {
	respondNotImplemented(w, r)
}
