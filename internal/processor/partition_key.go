package processor

import (
	"time"

	"github.com/reportportal/service-ingest/internal/data/buffer"
)

type partitionKey struct {
	ProjectKey  string
	LaunchUUID  string
	EntityType  buffer.EntityType
	EventDate   string
	WindowStart int64
}

func newPartitionKey(event buffer.EventEnvelope, window time.Duration) partitionKey {
	return partitionKey{
		ProjectKey:  event.ProjectKey,
		LaunchUUID:  event.LaunchUUID,
		EntityType:  event.EntityType,
		EventDate:   event.Timestamp.Format("2006-01-02"),
		WindowStart: event.Timestamp.Truncate(window).Unix(),
	}
}
