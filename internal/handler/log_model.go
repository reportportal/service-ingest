package handler

import (
	"net/http"
	"time"

	"github.com/reportportal/service-ingest/internal/model"
)

type SaveLogBatchRQ []SaveLogRQ

func (rq *SaveLogBatchRQ) Bind(r *http.Request) error {
	for i := range *rq {
		if err := (*rq)[i].Bind(r); err != nil {
			return err
		}
	}

	return nil
}

type SaveLogRQ struct {
	UUID       string         `json:"uuid,omitempty" validate:"omitempty,uuid"`
	LaunchUUID string         `json:"launchUuid" validate:"required,uuid"`
	ItemUUID   string         `json:"itemUuid,omitempty" validate:"omitempty,uuid"`
	Timestamp  time.Time      `json:"time" validate:"required"`
	Level      model.LogLevel `json:"level,omitempty" validate:"omitempty,oneof=error warn info debug trace fatal unknown"`
	Message    string         `json:"message,omitempty"`
	File       LogFile        `json:"file,omitempty"`
}

func (rq *SaveLogRQ) Bind(_ *http.Request) error {
	if err := validate.Struct(rq); err != nil {
		return err
	}

	if rq.Level == "" {
		rq.Level = model.LogLevelUnknown
	}

	return nil
}

func (rq *SaveLogRQ) toLogModel() model.Log {
	return model.Log{
		UUID:       rq.UUID,
		ItemUUID:   rq.ItemUUID,
		LaunchUUID: rq.LaunchUUID,
		Timestamp:  rq.Timestamp,
		Level:      rq.Level,
		Message:    rq.Message,
		File:       model.LogFile(rq.File),
	}
}

type LogFile struct {
	Name string `json:"name,omitempty"`
}

type SaveLogRS struct {
	Responses []LogResponse `json:"responses"`
}

func (rs *SaveLogRS) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type LogResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	StackTrace string `json:"stackTrace"`
}
