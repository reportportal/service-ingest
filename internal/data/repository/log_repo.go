package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/catalog"
	"github.com/reportportal/service-ingest/internal/model"
)

type LogRepositoryImpl struct {
	buffer  buffer.Buffer
	staging buffer.FileBuffer
}

func NewLogRepository(buffer buffer.Buffer, staging buffer.FileBuffer) *LogRepositoryImpl {
	return &LogRepositoryImpl{buffer, staging}
}

func (l *LogRepositoryImpl) Create(project string, log model.Log) error {
	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	envelope := buffer.EventEnvelope{
		ID:         uuid.New().String(),
		ProjectKey: project,
		LaunchUUID: log.LaunchUUID,
		EntityUUID: log.UUID,
		EntityType: buffer.EntityTypeLog,
		Operation:  buffer.OperationTypeCreate,
		Timestamp:  time.Now().UTC(),
		Data:       data,
		Size:       int64(len(data)),
	}

	if err := l.buffer.Put(context.Background(), envelope); err != nil {
		return fmt.Errorf("failed to put log create event: %w", err)
	}

	return nil
}

func (l *LogRepositoryImpl) CreateLogs(project string, logs []model.Log, files []*multipart.FileHeader) error {
	dic := make(map[string]*multipart.FileHeader)

	for _, f := range files {
		dic[f.Filename] = f
	}

	for _, log := range logs {
		if f, ok := dic[log.File.Name]; ok {
			hash, err := l.staging.Save(catalog.BuildFilePath(project, log.LaunchUUID), f)
			if err != nil {
				return fmt.Errorf("failed to save file %s: %w", f.Filename, err)
			}
			log.File.Hash = hash
		}

		if err := l.Create(project, log); err != nil {
			return fmt.Errorf("failed to create log %s: %w", log.UUID, err)
		}
	}

	return nil
}
