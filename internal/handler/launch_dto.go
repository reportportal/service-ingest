package handler

import (
	"time"

	"github.com/reportportal/service-ingest/internal/model"
)

type StartLaunchRQ struct {
	UUID        string           `json:"uuid" validate:"required,uuid"`
	Name        string           `json:"name" validate:"required"`
	StartTime   time.Time        `json:"startTime" validate:"required"`
	Description string           `json:"description,omitempty"`
	Attributes  []ItemAttribute  `json:"attributes,omitempty" validate:"omitempty,max=256,dive"`
	Mode        model.LaunchMode `json:"mode,omitempty" validate:"omitempty,oneof=DEFAULT DEBUG"`
	IsRerun     bool             `json:"rerun,omitempty"`
	RerunOf     string           `json:"rerunOf,omitempty" validate:"omitempty,uuid"`
}

type StartLaunchRS struct {
	ID     string `json:"id"`
	Number int64  `json:"number"`
}

type FinishLaunchRQ struct {
	EndTime     time.Time          `json:"endTime" validate:"required"`
	Status      model.LaunchStatus `json:"status,omitempty" validate:"omitempty,oneof=PASSED FAILED STOPPED SKIPPED INTERRUPTED CANCELLED INFO WARN"`
	Description string             `json:"description,omitempty"`
	Attributes  []ItemAttribute    `json:"attributes,omitempty" validate:"omitempty,max=256,dive"`
}

type FinishLaunchRS struct {
	ID     string `json:"id"`
	Number int64  `json:"number"`
	Link   string `json:"link"`
}

// GetLaunchRS represents launch resource with dates in ISO-8601 format
type GetLaunchRS struct {
	ID                  int64                  `json:"id" validate:"required"`
	UUID                string                 `json:"uuid" validate:"required"`
	Name                string                 `json:"name" validate:"required"`
	Number              int64                  `json:"number" validate:"required"`
	Description         string                 `json:"description,omitempty"`
	Status              string                 `json:"status" validate:"required"`
	Owner               string                 `json:"owner"`
	StartTime           time.Time              `json:"startTime" validate:"required"`
	EndTime             time.Time              `json:"endTime,omitempty"`
	LastModified        time.Time              `json:"lastModified"`
	Mode                model.LaunchMode       `json:"mode,omitempty"`
	Statistics          Statistics             `json:"statistics,omitempty"`
	Attributes          []ItemAttribute        `json:"attributes,omitempty"`
	Analysing           []string               `json:"analysing,omitempty"`
	ApproximateDuration float64                `json:"approximateDuration,omitempty"`
	HasRetries          bool                   `json:"hasRetries,omitempty"`
	Rerun               bool                   `json:"rerun,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	RetentionPolicy     string                 `json:"retentionPolicy,omitempty" validate:"omitempty,oneof=IMPORTANT REGULAR"`
}

// GetLaunchOldRS represents legacy launch resource with dates in Unix timestamp format (milliseconds since epoch)
// StartTime - Unix timestamp in milliseconds
// EndTime - Unix timestamp in milliseconds
// LastModified - Unix timestamp in milliseconds
type GetLaunchOldRS struct {
	ID                  int64                  `json:"id" validate:"required"`
	UUID                string                 `json:"uuid" validate:"required"`
	Name                string                 `json:"name" validate:"required"`
	Number              int64                  `json:"number" validate:"required"`
	Description         string                 `json:"description,omitempty"`
	Status              string                 `json:"status" validate:"required"`
	Owner               string                 `json:"owner"`
	StartTime           int64                  `json:"startTime" validate:"required"`
	EndTime             int64                  `json:"endTime,omitempty"`
	LastModified        int64                  `json:"lastModified"`
	Mode                model.LaunchMode       `json:"mode,omitempty"`
	Statistics          Statistics             `json:"statistics,omitempty"`
	Attributes          []ItemAttribute        `json:"attributes,omitempty"`
	Analysing           []string               `json:"analysing,omitempty"`
	ApproximateDuration float64                `json:"approximateDuration,omitempty"`
	HasRetries          bool                   `json:"hasRetries,omitempty"`
	Rerun               bool                   `json:"rerun,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	RetentionPolicy     string                 `json:"retentionPolicy,omitempty" validate:"omitempty,oneof=IMPORTANT REGULAR"`
}

type UpdateLaunchRQ struct {
	Description string           `json:"description,omitempty"`
	Mode        model.LaunchMode `json:"mode,omitempty" validate:"omitempty,oneof=DEFAULT DEBUG"`
	Attributes  []ItemAttribute  `json:"attributes,omitempty" validate:"omitempty,max=256,dive"`
}

type UpdateLaunchRS struct {
	Message string `json:"message"`
}

type Statistics struct {
	Executions map[string]int64            `json:"executions,omitempty"`
	Defects    map[string]map[string]int32 `json:"defects,omitempty"`
}

func (sl StartLaunchRQ) toLaunchModel() model.Launch {
	return model.Launch{
		ID:          sl.UUID,
		UUID:        sl.UUID,
		Number:      0,
		Name:        sl.Name,
		Description: sl.Description,
		StartTime:   sl.StartTime,
		Mode:        sl.Mode,
		Attributes:  sl.toAttributesModel(),
		IsRerun:     sl.IsRerun,
		RerunOf:     sl.RerunOf,
	}
}

func (sl StartLaunchRQ) toAttributesModel() []model.Attribute {
	attrs := make([]model.Attribute, 0, len(sl.Attributes))
	for _, a := range sl.Attributes {
		attrs = append(attrs, a.toModelAttribute())
	}
	return attrs
}
