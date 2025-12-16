package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/reportportal/service-ingest/internal/model"
	"github.com/reportportal/service-ingest/internal/service"
)

type LogHandler struct {
	logService *service.LogService
}

func NewLogHandler(logService *service.LogService) *LogHandler {
	return &LogHandler{logService: logService}
}

func (h LogHandler) routes() chi.Router {
	r := chi.NewRouter()

	r.With(ParseMultipartForm).Post("/", h.createLog)

	return r
}

func (h LogHandler) createLog(w http.ResponseWriter, r *http.Request) {
	if r.MultipartForm != nil {
		h.saveBatch(w, r)
		return
	}

	h.saveLog(w, r)
}

func (h LogHandler) saveBatch(w http.ResponseWriter, r *http.Request) {
	defer r.MultipartForm.RemoveAll()

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

func (h LogHandler) saveLog(w http.ResponseWriter, r *http.Request) {
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
