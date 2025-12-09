package model

import "time"

type Launch struct {
	ID          int64        `json:"id"`
	UUID        string       `json:"uuid"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Status      LaunchStatus `json:"status"`
	Owner       string       `json:"owner"`
	StartTime   time.Time    `json:"stat_time"`
	EndTime     *time.Time   `json:"end_time"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Mode        LaunchMode   `json:"mode"`
	Statistics  Statistics   `json:"statistics"`
	Attributes  []Attribute  `json:"attributes"`
	IsRerun     bool         `json:"isRerun"`
	RerunOf     string       `json:"rerunOf"`
	HasRetries  bool         `json:"hasRetries"`
	// Number      int64        `json:"number" `
}

type Statistics struct {
	Executions map[string]int64            `json:"executions"`
	Defects    map[string]map[string]int32 `json:"defects"`
}

// Duration returns the duration of the launch in seconds.
// If the launch is still in progress, it returns the duration from the start time to the current time.
func (l *Launch) Duration() float64 {
	if l.EndTime == nil {
		return time.Since(l.StartTime).Seconds()
	}

	return l.EndTime.Sub(l.StartTime).Seconds()
}
