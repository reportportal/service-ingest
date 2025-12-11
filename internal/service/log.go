package service

import (
	"mime/multipart"

	"github.com/reportportal/service-ingest/internal/model"
)

type LogService struct {
	logRepo LogRepository
}

func NewLogService(repo LogRepository) *LogService {
	return &LogService{repo}
}

func (s *LogService) SaveLog(project string, log model.Log) error {
	if err := s.logRepo.Create(project, log); err != nil {
		return err
	}

	return nil
}

func (s *LogService) SaveLogs(project string, logs []model.Log, files []*multipart.FileHeader) error {
	if err := s.logRepo.CreateLogs(project, logs, files); err != nil {
		return err
	}

	return nil
}
