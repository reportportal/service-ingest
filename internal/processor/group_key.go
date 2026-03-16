package processor

import (
	"github.com/reportportal/service-ingest/internal/data/buffer"
)

type groupKey struct {
	ProjectKey string
	LaunchUUID string
	EntityType buffer.EntityType
	EventDate  string
}

func newGroupKey(event buffer.EventEnvelope) groupKey {
	return groupKey{
		ProjectKey: event.ProjectKey,
		LaunchUUID: event.LaunchUUID,
		EntityType: event.EntityType,
		EventDate:  event.Timestamp.Format("2006-01-02"),
	}
}
