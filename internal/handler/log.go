package handler

import (
	"net/http"
	"strings"

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

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		h.saveBatch(w, r)
		return
	}

	h.saveLog(w, r)
}

func (h logHandler) saveBatch(w http.ResponseWriter, r *http.Request) {
	batch := &SaveLogBatchRQ{}
	if err := batch.MultipartFormBind(r); err != nil {
		render.Render(w, r, InvalidRequestError(err))
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

	rs := make([]LogResponse, len(logs))
	for i, log := range logs {
		rs[i] = LogResponse{ID: log.UUID}
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, &SaveLogRS{Responses: rs})
}

func (h logHandler) saveLog(w http.ResponseWriter, r *http.Request) {
	data := &SaveLogRQ{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, InvalidRequestError(err))
		return
	}

	log := data.toLogModel()
	projectName := chi.URLParam(r, "projectName")
	if err := h.logService.SaveLog(projectName, log); err != nil {
		render.Render(w, r, InternalServerError)
		return
	}

	rs := []LogResponse{{ID: log.UUID}}
	render.Status(r, http.StatusOK)
	render.Render(w, r, &SaveLogRS{Responses: rs})
}
