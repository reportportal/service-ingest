package buffer

import (
	"encoding/json"
	"time"
)

type EventEnvelope struct {
	ID         string          `json:"id"`
	BufferKey  []byte          `json:"buffer_key"`
	ProjectKey string          `json:"project_key"`
	LaunchUUID string          `json:"launch_uuid"`
	EntityUUID string          `json:"entity_uuid"`
	EntityType EntityType      `json:"entity_type"`
	Operation  OperationType   `json:"operation"`
	Timestamp  time.Time       `json:"timestamp"`
	Data       json.RawMessage `json:"data"`
	Size       int64           `json:"size"`
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
