package service

import (
	"mime/multipart"

	"github.com/reportportal/service-ingest/internal/model"
)

type LogRepository interface {
	Create(project string, logs model.Log) error
	CreateLogs(project string, logs []model.Log, files []*multipart.FileHeader) error
}
