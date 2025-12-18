package repository

import (
	"mime/multipart"

	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type LogRepositoryImpl struct {
	buffer buffer.Buffer
}

func NewLogRepository(buffer buffer.Buffer) *LogRepositoryImpl {
	return &LogRepositoryImpl{buffer}
}

func (l LogRepositoryImpl) Create(project string, logs model.Log) error {
	//TODO implement me
	panic("implement me")
}

func (l LogRepositoryImpl) CreateLogs(project string, logs []model.Log, files []*multipart.FileHeader) error {
	//TODO implement me
	panic("implement me")
}
