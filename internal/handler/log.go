package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/reportportal/service-ingest/internal/model"
	"github.com/reportportal/service-ingest/internal/service"
)

type logHandler struct {
	logService service.LogService
}

func (h logHandler) routes() chi.Router {
	r := chi.NewRouter()

	r.With(ParseMultipartForm).Post("/", h.createLog)

	return r
}

func (h logHandler) createLog(w http.ResponseWriter, r *http.Request) {
	batch := &SaveLogBatchRQ{}

	if err := render.Bind(r, batch); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	logs := make([]model.Log, len(*batch))
	for i, log := range *batch {
		logs[i] = log.toLogModel()
	}
	files := r.MultipartForm.File["file"]
	projectName := chi.URLParam(r, "projectName")

	if err := h.logService.SaveLogs(projectName, logs, files); err != nil {
		render.Render(w, r, InternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &SaveLogRS{
		Responses: []LogResponse{
			{
				ID:         "",
				Message:    "Logs saved successfully",
				StackTrace: "",
			},
		},
	})
}
