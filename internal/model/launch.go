package model

import "time"

const (
	LaunchModeDefault LaunchMode = "DEFAULT"
	LaunchModeDebug   LaunchMode = "DEBUG"
)

const (
	LaunchStatusPassed      LaunchStatus = "PASSED"
	LaunchStatusFailed      LaunchStatus = "FAILED"
	LaunchStatusStopped     LaunchStatus = "STOPPED"
	LaunchStatusSkipped     LaunchStatus = "SKIPPED"
	LaunchStatusInterrupted LaunchStatus = "INTERRUPTED"
	LaunchStatusCancelled   LaunchStatus = "CANCELLED"
	LaunchStatusInfo        LaunchStatus = "INFO"
	LaunchStatusWarn        LaunchStatus = "WARN"
)

type LaunchMode string

type LaunchStatus string

type Launch struct {
	ID          string       `json:"id"`
	UUID        string       `json:"uuid"`
	Number      int64        `json:"number"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Status      LaunchStatus `json:"status"`
	Owner       string       `json:"owner"`
	StartTime   time.Time    `json:"stat_time"`
	EndTime     *time.Time   `json:"end_time,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Mode        LaunchMode   `json:"mode,omitempty"`
	Attributes  []Attribute  `json:"attributes,omitempty"`
	IsRerun     bool         `json:"isRerun,omitempty"`
	RerunOf     string       `json:"rerunOf,omitempty"`
}
