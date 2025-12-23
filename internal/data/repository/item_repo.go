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

type ItemRepositoryImpl struct {
	buffer buffer.Buffer
}

func NewItemRepository(buffer buffer.Buffer) *ItemRepositoryImpl {
	return &ItemRepositoryImpl{buffer}
}

func (i *ItemRepositoryImpl) Get(project string, itemUUID string) (*model.Item, error) {
	// TODO: implement query side (read from parquet/materialized view)
	return nil, fmt.Errorf("not implemented")
}

func (i *ItemRepositoryImpl) Start(project string, item model.Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	envelope := buffer.EventEnvelope{
		ID:         uuid.New().String(),
		ProjectKey: project,
		LaunchUUID: item.LaunchUUID,
		EntityUUID: item.UUID,
		EntityType: buffer.EntityTypeItem,
		Operation:  buffer.OperationTypeStart,
		Timestamp:  time.Now().UTC(),
		Data:       data,
		Size:       int64(len(data)),
	}

	if err := i.buffer.Put(context.Background(), envelope); err != nil {
		return fmt.Errorf("failed to put item start event: %w", err)
	}

	return nil
}

func (i *ItemRepositoryImpl) Finish(project string, item model.Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	envelope := buffer.EventEnvelope{
		ID:         uuid.New().String(),
		ProjectKey: project,
		LaunchUUID: item.LaunchUUID,
		EntityUUID: item.UUID,
		EntityType: buffer.EntityTypeItem,
		Operation:  buffer.OperationTypeFinish,
		Timestamp:  time.Now().UTC(),
		Data:       data,
		Size:       int64(len(data)),
	}

	if err := i.buffer.Put(context.Background(), envelope); err != nil {
		return fmt.Errorf("failed to put item finish event: %w", err)
	}

	return nil
}
