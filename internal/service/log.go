package service

import (
	"fmt"
	"mime/multipart"

	"github.com/reportportal/service-ingest/internal/model"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) SaveLog(project string, log model.Log) error {
	return nil
}

func (s *LogService) SaveLogs(project string, logs []model.Log, files []*multipart.FileHeader) error {
	fmt.Printf("Saving %d logs for project %s\n", len(logs), project)
	fmt.Printf("Saving %d files for project %s\n", len(files), project)
	return nil
}
