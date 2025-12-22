package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type LaunchRepositoryImpl struct {
	buffer buffer.Buffer
}

func NewLaunchRepo(buffer buffer.Buffer) *LaunchRepositoryImpl {
	return &LaunchRepositoryImpl{buffer}
}

func (l *LaunchRepositoryImpl) Get(project string, launchUUID string) (*model.Launch, error) {
	// TODO: implement query side (read from parquet/materialized view)
	return nil, fmt.Errorf("not implemented")
}

func (l *LaunchRepositoryImpl) Start(project string, launch model.Launch) error {
	data, err := json.Marshal(launch)
	if err != nil {
		return fmt.Errorf("failed to marshal launch: %w", err)
	}

	envelope := buffer.EventEnvelope{
		ID:         uuid.New().String(),
		Project:    project,
		EntityUUID: launch.UUID,
		EntityType: buffer.EntityTypeLaunch,
		Operation:  buffer.OperationTypeStart,
		Timestamp:  time.Now().UTC(),
		Data:       data,
		Size:       int64(len(data)),
	}

	if err := l.buffer.Put(context.Background(), envelope); err != nil {
		return fmt.Errorf("failed to put launch start event: %w", err)
	}

	return nil
}

func (l *LaunchRepositoryImpl) Update(project string, launch model.Launch) error {
	data, err := json.Marshal(launch)
	if err != nil {
		return fmt.Errorf("failed to marshal launch: %w", err)
	}

	envelope := buffer.EventEnvelope{
		ID:         uuid.New().String(),
		Project:    project,
		EntityUUID: launch.UUID,
		EntityType: buffer.EntityTypeLaunch,
		Operation:  buffer.OperationTypeUpdate,
		Timestamp:  time.Now().UTC(),
		Data:       data,
		Size:       int64(len(data)),
	}

	if err := l.buffer.Put(context.Background(), envelope); err != nil {
		return fmt.Errorf("failed to put launch update event: %w", err)
	}

	return nil
}

func (l *LaunchRepositoryImpl) Finish(project string, launch model.Launch) error {
	data, err := json.Marshal(launch)
	if err != nil {
		return fmt.Errorf("failed to marshal launch: %w", err)
	}

	envelope := buffer.EventEnvelope{
		ID:         uuid.New().String(),
		Project:    project,
		EntityUUID: launch.UUID,
		EntityType: buffer.EntityTypeLaunch,
		Operation:  buffer.OperationTypeFinish,
		Timestamp:  time.Now().UTC(),
		Data:       data,
		Size:       int64(len(data)),
	}

	if err := l.buffer.Put(context.Background(), envelope); err != nil {
		return fmt.Errorf("failed to put launch finish event: %w", err)
	}

	return nil
}
