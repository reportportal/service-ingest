package scheme

import (
	"time"

	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type LaunchEvent struct {
	ID          string               `parquet:"id"`
	Timestamp   time.Time            `parquet:"timestamp"`
	Operation   buffer.OperationType `parquet:"operation"`
	UUID        string               `parquet:"uuid"`
	Name        string               `parquet:"name"`
	Description string               `parquet:"description"`
	Status      model.LaunchStatus   `parquet:"status"`
	Owner       string               `parquet:"owner"`
	StartTime   time.Time            `parquet:"start_time"`
	EndTime     *time.Time           `parquet:"end_time,optional"`
	UpdatedAt   time.Time            `parquet:"updated_at"`
	Mode        model.LaunchMode     `parquet:"mode"`
	Attributes  []model.Attribute    `parquet:"attributes"`
	IsRerun     bool                 `parquet:"is_rerun"`
	RerunOf     string               `parquet:"rerun_of"`
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
