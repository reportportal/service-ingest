package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type LogRepositoryImpl struct {
	buffer buffer.Buffer
}

func NewLogRepository(buffer buffer.Buffer) *LogRepositoryImpl {
	return &LogRepositoryImpl{buffer}
}

func (l *LogRepositoryImpl) Create(project string, log model.Log) error {
	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	envelope := buffer.EventEnvelope{
		ID:         uuid.New().String(),
		Project:    project,
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

func (l *LogRepositoryImpl) CreateLogs(project string, logs []model.Log, _ []*multipart.FileHeader) error {
	// TODO: handle file uploads (store in object storage, add reference to log)

	for _, log := range logs {
		if err := l.Create(project, log); err != nil {
			return fmt.Errorf("failed to create log %s: %w", log.UUID, err)
		}
	}

	return nil
}
