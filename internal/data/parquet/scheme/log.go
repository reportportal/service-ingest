package scheme

import (
	"time"

	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type LogEvent struct {
	ID         string               `parquet:"id"`
	Timestamp  time.Time            `parquet:"timestamp"`
	Operation  buffer.OperationType `parquet:"operation"`
	UUID       string               `parquet:"uuid"`
	ItemUUID   string               `parquet:"item_uuid"`
	LaunchUUID string               `parquet:"launch_uuid"`
	LogTime    time.Time            `parquet:"log_time"`
	Level      model.LogLevel       `parquet:"level"`
	Message    string               `parquet:"message"`
	File       model.LogFile        `parquet:"file"`
}

func NewLogEvent(event buffer.EventEnvelope, log model.Log) LogEvent {
	return LogEvent{
		ID:         event.ID,
		Timestamp:  event.Timestamp,
		Operation:  event.Operation,
		UUID:       log.UUID,
		ItemUUID:   log.ItemUUID,
		LaunchUUID: log.LaunchUUID,
		LogTime:    log.Timestamp,
		Level:      log.Level,
		Message:    log.Message,
	}
}
