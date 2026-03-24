package scheme

import (
	"time"

	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type LaunchEvent struct {
	ID          string               `parquet:"event_id,delta"`
	Timestamp   time.Time            `parquet:"event_timestamp,timestamp"`
	Operation   buffer.OperationType `parquet:"operation,enum"`
	UUID        string               `parquet:"uuid,delta"`
	Name        string               `parquet:"name,optional,dict"`
	Description string               `parquet:"description,optional,string"`
	Status      model.LaunchStatus   `parquet:"status,enum"`
	Owner       string               `parquet:"owner,dict"`
	StartTime   *time.Time           `parquet:"start_time,optional,timestamp"`
	EndTime     *time.Time           `parquet:"end_time,optional,timestamp"`
	UpdatedAt   time.Time            `parquet:"updated_at,timestamp"`
	Mode        model.LaunchMode     `parquet:"mode,optional,enum"`
	Attributes  []model.Attribute    `parquet:"attributes,list"`
	IsRerun     bool                 `parquet:"is_rerun"`
	RerunOf     string               `parquet:"rerun_of,optional,delta"`
}

func NewLaunchEvent(event buffer.EventEnvelope, launch model.Launch) LaunchEvent {
	return LaunchEvent{
		ID:          event.ID,
		Timestamp:   event.Timestamp,
		Operation:   event.Operation,
		UUID:        launch.UUID,
		Name:        launch.Name,
		Description: launch.Description,
		Status:      launch.Status,
		Owner:       launch.Owner,
		StartTime:   launch.StartTime,
		EndTime:     launch.EndTime,
		UpdatedAt:   launch.UpdatedAt,
		Mode:        launch.Mode,
		Attributes:  launch.Attributes,
		IsRerun:     launch.IsRerun,
		RerunOf:     launch.RerunOf,
	}
}
