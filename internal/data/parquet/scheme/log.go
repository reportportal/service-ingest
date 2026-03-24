package scheme

import (
	"time"

	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type LogEvent struct {
	ID         string               `parquet:"event_id,delta"`
	Timestamp  time.Time            `parquet:"event_timestamp,timestamp"`
	Operation  buffer.OperationType `parquet:"operation,enum"`
	UUID       string               `parquet:"uuid,delta"`
	ItemUUID   string               `parquet:"item_uuid,delta"`
	LaunchUUID string               `parquet:"launch_uuid,delta"`
	LogTime    time.Time            `parquet:"log_time,timestamp"`
	Level      model.LogLevel       `parquet:"level,enum"`
	Message    string               `parquet:"message,string"`
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
		File:       log.File,
	}
}
