package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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
	data := &StartLaunchRQ{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	projectName := chi.URLParam(r, "projectName")
	err := h.service.StartLaunch(projectName, data.toLaunchModel())
	if err != nil {
		render.Render(w, r, InternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &StartLaunchRS{UUID: data.UUID})
}

func (h launchHandler) finishLaunch(w http.ResponseWriter, r *http.Request) {
	data := &FinishLaunchRQ{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	launchUUID := chi.URLParam(r, "launchUuid")
	projectName := chi.URLParam(r, "projectName")

	err := h.service.FinishLaunch(projectName, launchUUID, data.toFinishLaunchModel())
	if err != nil {
		render.Render(w, r, InternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &FinishLaunchRS{
		UUID: launchUUID,
		Link: r.Host + "/ui/#/" + projectName + "/launches/" + launchUUID,
	})
}

func (h launchHandler) updateLaunch(w http.ResponseWriter, r *http.Request) {
	data := &UpdateLaunchRQ{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	projectName := chi.URLParam(r, "projectName")
	launchId, err := strconv.ParseInt(chi.URLParam(r, "launchId"), 10, 64)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("invalid launch ID format")))
		return
	}

	err = h.service.UpdateLaunch(projectName, launchId, data.toUpdateLaunchModel())
	if err != nil {
		render.Render(w, r, InternalServerError)
		return
	}

	//render.Status(r, http.StatusOK)
	//render.Render(w, r, &UpdateLaunchRS{Message: "Launch updated successfully"})
	RespondNotImplemented(w, r)
}

func (h launchHandler) getLaunch(w http.ResponseWriter, r *http.Request) {
	projectName := chi.URLParam(r, "projectName")
	launchUUID := chi.URLParam(r, "launchUuid")

	_, err := h.service.GetLaunch(projectName, launchUUID)
	if err != nil {
		render.Render(w, r, InternalServerError)
		return
	}

	//render.Status(r, http.StatusOK)
	//render.Render(w, r, NewGetLaunchOldRS(launch))
	RespondNotImplemented(w, r)
}

func (h launchHandler) mergeLaunch(w http.ResponseWriter, r *http.Request) {
	RespondNotImplemented(w, r)
}
