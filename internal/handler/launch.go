package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/reportportal/service-ingest/internal/service"
)

type LaunchHandler struct {
	service *service.LaunchService
}

func NewLaunchHandler(svc *service.LaunchService) *LaunchHandler {
	return &LaunchHandler{service: svc}
}

func (h LaunchHandler) routesV1() chi.Router {
	r := chi.NewRouter()

	r.Get("/uuid/{launchUuid}", h.getLaunch)
	r.Put("/{launchId}/update", h.updateLaunch)

	return r
}

func (h LaunchHandler) routesV2() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.startLaunch)
	r.Post("/merge", h.mergeLaunch)
	r.Put("/{launchUuid}/finish", h.finishLaunch)

	return r
}

func (h LaunchHandler) startLaunch(w http.ResponseWriter, r *http.Request) {
	data := &StartLaunchRQ{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, InvalidRequestError(err))
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

func (h LaunchHandler) finishLaunch(w http.ResponseWriter, r *http.Request) {
	data := &FinishLaunchRQ{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, InvalidRequestError(err))
		return
	}

	launchUUID := chi.URLParam(r, "launchUuid")
	projectName := chi.URLParam(r, "projectName")

	err := h.service.FinishLaunch(projectName, launchUUID, data.toLaunchModel())
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

func (h LaunchHandler) updateLaunch(w http.ResponseWriter, r *http.Request) {
	RespondNotImplemented(w, r)
	//data := &UpdateLaunchRQ{}
	//
	//if err := render.Bind(r, data); err != nil {
	//	render.Render(w, r, InvalidRequestError(err))
	//	return
	//}
	//
	//projectName := chi.URLParam(r, "projectName")
	//launchId, err := strconv.ParseInt(chi.URLParam(r, "launchId"), 10, 64)
	//if err != nil {
	//	render.Render(w, r, InvalidRequestError(errors.New("invalid launch ID format")))
	//	return
	//}
	//
	//err = h.service.UpdateLaunch(projectName, launchId, data.toLaunchModel())
	//if err != nil {
	//	render.Render(w, r, InternalServerError)
	//	return
	//}

	//render.Status(r, http.StatusOK)
	//render.Render(w, r, &UpdateLaunchRS{Message: "Launch updated successfully"})
}

func (h LaunchHandler) getLaunch(w http.ResponseWriter, r *http.Request) {
	RespondNotImplemented(w, r)
	//projectName := chi.URLParam(r, "projectName")
	//launchUUID := chi.URLParam(r, "launchUuid")
	//
	//_, err := h.service.GetLaunch(projectName, launchUUID)
	//if err != nil {
	//	render.Render(w, r, InternalServerError)
	//	return
	//}
	//
	//render.Status(r, http.StatusOK)
	//render.Render(w, r, NewGetLaunchOldRS(launch))
}

func (h LaunchHandler) mergeLaunch(w http.ResponseWriter, r *http.Request) {
	RespondNotImplemented(w, r)
}
