package buffer

import (
	"encoding/json"
	"time"
)

type EventEnvelope struct {
	ID         string          `json:"id"`
	EntityUUID string          `json:"entity_uuid"`
	RequestID  string          `json:"request_id"`
	EntityType string          `json:"entity_type"`
	Operation  string          `json:"operation"`
	Timestamp  time.Time       `json:"timestamp"`
	Data       json.RawMessage `json:"Data"`
	Size       int64           `json:"size"`
}
