package scheme

import (
	"time"

	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type ItemEvent struct {
	ID          string               `parquet:"id"`
	Timestamp   time.Time            `parquet:"timestamp"`
	Operation   buffer.OperationType `parquet:"operation"`
	UUID        string               `parquet:"uuid"`
	LaunchUUID  string               `parquet:"launch_uuid"`
	Name        string               `parquet:"name"`
	Description string               `parquet:"description"`
	Type        model.ItemType       `parquet:"type"`
	Status      model.ItemStatus     `parquet:"status"`
	StartTime   time.Time            `parquet:"start_time"`
	EndTime     *time.Time           `parquet:"end_time,optional"`
	UpdatedAt   time.Time            `parquet:"updated_at"`
	Attributes  model.Attributes     `parquet:"attributes"`
	Parameters  model.Parameters     `parquet:"parameters"`
	CodeRef     string               `parquet:"code_ref"`
	TestCaseId  string               `parquet:"test_case_id"`
	ParentUUID  string               `parquet:"parent_uuid"`
	IsRetry     bool                 `parquet:"is_retry"`
	RetryOf     string               `parquet:"retry_of"`
	Issue       model.Issue          `parquet:"issue"`
}

func NewItemEvent(event buffer.EventEnvelope, item model.Item) ItemEvent {
	return ItemEvent{
		ID:         event.ID,
		Timestamp:  event.Timestamp,
		Operation:  event.Operation,
		UUID:       item.UUID,
		LaunchUUID: item.LaunchUUID,
		Name:       item.Name,
		Type:       item.Type,
		Status:     item.Status,
		StartTime:  item.StartTime,
		ParentUUID: item.ParentUUID,
		IsRetry:    item.IsRetry,
	}
}
