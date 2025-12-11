package model

import "time"

type Launch struct {
	UUID        string
	Name        string
	Description string
	Status      LaunchStatus
	Owner       string
	StartTime   time.Time
	EndTime     *time.Time
	UpdatedAt   time.Time
	Mode        LaunchMode
	Statistics  Statistics
	Attributes  []Attribute
	IsRerun     bool
	RerunOf     string
	HasRetries  bool
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
