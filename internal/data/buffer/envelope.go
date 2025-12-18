package buffer

import (
	"encoding/json"
	"time"
)

type EventEnvelope struct {
	ID         string          `json:"id"`
	EntityUUID string          `json:"entity_uuid"`
	EntityType string          `json:"entity_type"`
	Operation  string          `json:"operation"`
	Timestamp  time.Time       `json:"timestamp"`
	Data       json.RawMessage `json:"data"`
	Size       int64           `json:"size"`

	//Tracing
	RequestID string `json:"request_id,omitempty"`

	//Lease management
	LeaseID        string     `json:"lease_id,omitempty"`
	LeaseExpiresAt *time.Time `json:"lease_expires_at,omitempty"`
}

func (e *EventEnvelope) IsAvailable() bool {
	if e.LeaseID == "" {
		return true
	}
	if e.LeaseExpiresAt != nil && time.Now().After(*e.LeaseExpiresAt) {
		return true
	}
	return false
}
