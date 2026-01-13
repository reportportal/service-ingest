package handler

import (
	"net/http"
	"time"

	"github.com/reportportal/service-ingest/internal/model"
)

type StartLaunchRQ struct {
	UUID        string           `json:"uuid" validate:"omitempty,uuid"`
	Name        string           `json:"name" validate:"required"`
	StartTime   time.Time        `json:"startTime" validate:"required"`
	Description string           `json:"description,omitempty"`
	Attributes  Attributes       `json:"attributes,omitempty" validate:"omitempty,max=256,dive"`
	Mode        model.LaunchMode `json:"mode,omitempty" validate:"omitempty,oneof=DEFAULT DEBUG"`
	IsRerun     bool             `json:"rerun,omitempty"`
	RerunOf     string           `json:"rerunOf,omitempty" validate:"omitempty,uuid"`
}

func (rq *StartLaunchRQ) Bind(_ *http.Request) error {
	if err := validate.Struct(rq); err != nil {
		return err
	}

	if rq.Mode == "" {
		rq.Mode = model.LaunchModeDefault
	}

	return nil
}

func (rq *StartLaunchRQ) toLaunchModel() model.Launch {
	return model.Launch{
		UUID:        rq.UUID,
		Name:        rq.Name,
		Description: rq.Description,
		StartTime:   rq.StartTime,
		Mode:        rq.Mode,
		Attributes:  rq.Attributes.toAttributesModel(),
		IsRerun:     rq.IsRerun,
		RerunOf:     rq.RerunOf,
	}
}

type StartLaunchRS struct {
	UUID string `json:"id"`
}

func (rs *StartLaunchRS) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type FinishLaunchRQ struct {
	EndTime     time.Time          `json:"endTime" validate:"required"`
	Status      model.LaunchStatus `json:"status,omitempty" validate:"omitempty,oneof=PASSED FAILED STOPPED SKIPPED INTERRUPTED CANCELLED INFO WARN"`
	Description string             `json:"description,omitempty"`
	Attributes  Attributes         `json:"attributes,omitempty" validate:"omitempty,max=256,dive"`
}

func (rq *FinishLaunchRQ) Bind(_ *http.Request) error {
	if err := validate.Struct(rq); err != nil {
		return err
	}

	return nil
}

func (rq *FinishLaunchRQ) toLaunchModel() model.Launch {
	return model.Launch{
		EndTime:     &rq.EndTime,
		Status:      rq.Status,
		Description: rq.Description,
		Attributes:  rq.Attributes.toAttributesModel(),
	}
}

type FinishLaunchRS struct {
	UUID string `json:"id"`
	Link string `json:"link"`
}

func (rs *FinishLaunchRS) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type UpdateLaunchRQ struct {
	Description string           `json:"description,omitempty"`
	Mode        model.LaunchMode `json:"mode,omitempty" validate:"omitempty,oneof=DEFAULT DEBUG"`
	Attributes  Attributes       `json:"attributes,omitempty" validate:"omitempty,max=256,dive"`
}

func (rq *UpdateLaunchRQ) Bind(_ *http.Request) error {
	if err := validate.Struct(rq); err != nil {
		return err
	}

	return nil
}

func (rq *UpdateLaunchRQ) toLaunchModel() model.Launch {
	return model.Launch{
		Description: rq.Description,
		Mode:        rq.Mode,
		Attributes:  rq.Attributes.toAttributesModel(),
	}
}

type UpdateLaunchRS struct {
	Message string `json:"message"`
}

func (rs *UpdateLaunchRS) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// GetLaunchRS represents launch resource with dates in ISO-8601 format
type GetLaunchRS struct {
	StartTime    time.Time `json:"startTime" validate:"required"`
	EndTime      time.Time `json:"endTime,omitempty"`
	LastModified time.Time `json:"lastModified"`
	LaunchResource
}

func (rs *GetLaunchRS) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// GetLaunchOldRS represents legacy launch resource with dates in Unix timestamp format (milliseconds since epoch)
//
// StartTime, EndTime, LastModified - Unix timestamp in milliseconds
type GetLaunchOldRS struct {
	StartTime    int64 `json:"startTime" validate:"required"`
	EndTime      int64 `json:"endTime,omitempty"`
	LastModified int64 `json:"lastModified"`
	LaunchResource
}

func (rs *GetLaunchOldRS) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewGetLaunchOldRS(launch model.Launch) *GetLaunchOldRS {
	if launch.EndTime == nil {
		empty := time.Time{}
		launch.EndTime = &empty
	}

	return &GetLaunchOldRS{
		StartTime:    launch.StartTime.UnixMilli(),
		EndTime:      launch.EndTime.UnixMilli(),
		LastModified: launch.UpdatedAt.UnixMilli(),
		LaunchResource: LaunchResource{
			ID:                  0,
			UUID:                launch.UUID,
			Name:                launch.Name,
			Description:         launch.Description,
			Status:              launch.Status,
			Owner:               launch.Owner,
			ApproximateDuration: launch.Duration(),
			Mode:                launch.Mode,
			Statistics:          Statistics(launch.Statistics),
			Attributes:          fromAttributesModel(launch.Attributes),
			Rerun:               launch.IsRerun,
			HasRetries:          launch.HasRetries,
			Number:              0,
			Analysing:           []string{},
			Metadata:            map[string]interface{}{},
			RetentionPolicy:     "REGULAR",
		},
	}
}

type LaunchResource struct {
	ID                  int64                  `json:"id" validate:"required"`
	UUID                string                 `json:"uuid" validate:"required"`
	Name                string                 `json:"name" validate:"required"`
	Number              int64                  `json:"number" validate:"required"`
	Description         string                 `json:"description,omitempty"`
	Status              model.LaunchStatus     `json:"status" validate:"required"`
	Owner               string                 `json:"owner,omitempty"`
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

type Statistics struct {
	Executions map[string]int64            `json:"executions,omitempty"`
	Defects    map[string]map[string]int32 `json:"defects,omitempty"`
}
