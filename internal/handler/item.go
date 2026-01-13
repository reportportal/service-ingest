package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/reportportal/service-ingest/internal/service"
)

type ItemHandler struct {
	service *service.ItemService
}

func NewItemHandler(svc *service.ItemService) *ItemHandler {
	return &ItemHandler{service: svc}
}

func (h ItemHandler) routesV1() chi.Router {
	r := chi.NewRouter()

	r.Get("/uuid/{itemUuid}", h.getTestItem)

	return r
}

func (h ItemHandler) routesV2() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.startItem)
	r.Post("/{itemUuid}", h.startItem)
	r.Put("/{itemUuid}", h.finishTestItem)

	return r
}

func (h ItemHandler) startItem(w http.ResponseWriter, r *http.Request) {
	data := &StartTestItemRQ{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, InvalidRequestError(err))
		return
	}

	projectName := chi.URLParam(r, "projectName")

	uuid, err := h.service.StartItem(projectName, data.toItemModel())
	if err != nil {
		render.Render(w, r, InternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &StartTestItemRS{UUID: uuid})
}

func (h ItemHandler) finishTestItem(w http.ResponseWriter, r *http.Request) {
	data := &FinishTestItemRQ{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, InvalidRequestError(err))
		return
	}

	projectName := chi.URLParam(r, "projectName")
	itemUUID := chi.URLParam(r, "itemUuid")

	if err := h.service.FinishItem(projectName, itemUUID, data.toItemModel()); err != nil {
		render.Render(w, r, InternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &FinishTestItemRS{Message: "Item finished successfully"})
}

func (h ItemHandler) getTestItem(w http.ResponseWriter, r *http.Request) {
	projectName := chi.URLParam(r, "projectName")
	itemUUID := chi.URLParam(r, "itemUuid")

	item, err := h.service.GetItem(projectName, itemUUID)
	if err != nil {
		render.Render(w, r, InternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, NewTestItemResourceOldRS(item))
}
