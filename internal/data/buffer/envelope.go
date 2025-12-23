package buffer

import (
	"encoding/json"
	"time"
)

type EventEnvelope struct {
	ID         string          `json:"id"`
	ProjectKey string          `json:"project_key"`
	LaunchUUID string          `json:"launch_uuid"`
	EntityUUID string          `json:"entity_uuid"`
	EntityType EntityType      `json:"entity_type"`
	Operation  OperationType   `json:"operation"`
	Timestamp  time.Time       `json:"timestamp"`
	Data       json.RawMessage `json:"data"`
	Size       int64           `json:"size"`

	//Lease management
	LeaseID string `json:"lease_id,omitempty"`
	//LeaseExpiresAt *time.Time `json:"lease_expires_at,omitempty"`
}

func (e *EventEnvelope) IsAvailable() bool {
	if e.LeaseID == "" {
		return true
	}
	//if e.LeaseExpiresAt != nil && time.Now().After(*e.LeaseExpiresAt) {
	//	return true
	//}
	return false
}

type EntityType string

const (
	EntityTypeLaunch EntityType = "launch"
	EntityTypeItem   EntityType = "item"
	EntityTypeLog    EntityType = "log"
)

type OperationType string

const (
	OperationTypeStart  OperationType = "start"
	OperationTypeUpdate OperationType = "update"
	OperationTypeFinish OperationType = "finish"
	OperationTypeCreate OperationType = "create"
)
