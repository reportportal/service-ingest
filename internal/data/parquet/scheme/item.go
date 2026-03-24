package scheme

import (
	"time"

	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type ItemEvent struct {
	ID          string               `parquet:"event_id,delta"`
	Timestamp   time.Time            `parquet:"event_timestamp,timestamp"`
	Operation   buffer.OperationType `parquet:"operation,enum"`
	UUID        string               `parquet:"uuid,delta"`
	LaunchUUID  string               `parquet:"launch_uuid,delta"`
	Name        string               `parquet:"name,dict"`
	Description string               `parquet:"description,string"`
	Type        model.ItemType       `parquet:"type,enum"`
	Status      model.ItemStatus     `parquet:"status,enum"`
	StartTime   *time.Time           `parquet:"start_time,optional,timestamp"`
	EndTime     *time.Time           `parquet:"end_time,optional,timestamp"`
	UpdatedAt   time.Time            `parquet:"updated_at,timestamp"`
	Attributes  model.Attributes     `parquet:"attributes,list"`
	Parameters  model.Parameters     `parquet:"parameters,list"`
	CodeRef     string               `parquet:"code_ref,dict"`
	TestCaseId  string               `parquet:"test_case_id,delta"`
	ParentUUID  string               `parquet:"parent_uuid,delta"`
	IsRetry     bool                 `parquet:"is_retry"`
	RetryOf     string               `parquet:"retry_of,optional,delta"`
	Issue       model.Issue          `parquet:"issue"`
}

func NewItemEvent(event buffer.EventEnvelope, item model.Item) ItemEvent {
	return ItemEvent{
		ID:          event.ID,
		Timestamp:   event.Timestamp,
		Operation:   event.Operation,
		UUID:        item.UUID,
		LaunchUUID:  item.LaunchUUID,
		Name:        item.Name,
		Description: item.Description,
		Type:        item.Type,
		Status:      item.Status,
		StartTime:   item.StartTime,
		EndTime:     item.EndTime,
		UpdatedAt:   item.UpdatedAt,
		Attributes:  item.Attributes,
		Parameters:  item.Parameters,
		CodeRef:     item.CodeRef,
		TestCaseId:  item.TestCaseId,
		ParentUUID:  item.ParentUUID,
		IsRetry:     item.IsRetry,
		RetryOf:     item.RetryOf,
		Issue:       item.Issue,
	}
}
