package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/reportportal/service-ingest/internal/model"
)

type StartTestItemRQ struct {
	UUID        string         `json:"uuid" validate:"omitempty,uuid"`
	LaunchUUID  string         `json:"launchUuid" validate:"required,uuid"`
	StartTime   time.Time      `json:"startTime" validate:"required,datetime"`
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description,omitempty"`
	CodeRef     string         `json:"codeRef,omitempty"`
	TestCaseId  string         `json:"testCaseId,omitempty"`
	Attributes  Attributes     `json:"attributes,omitempty" validate:"max=256,dive,unique"`
	Parameters  Parameters     `json:"parameters,omitempty"`
	IsRetry     bool           `json:"retry,omitempty"`
	HasStats    bool           `json:"hasStats,omitempty"`
	RetryOf     string         `json:"retryOf,omitempty" validate:"omitempty,uuid"`
	Type        model.ItemType `json:"type" validate:"required,oneof=SUITE STORY TEST SCENARIO STEP BEFORE_CLASS BEFORE_GROUPS BEFORE_METHOD BEFORE_SUITE BEFORE_TEST AFTER_CLASS AFTER_GROUPS AFTER_METHOD AFTER_SUITE AFTER_TEST"`
	parentUUID  string
}

func (rq *StartTestItemRQ) Bind(r *http.Request) error {
	if err := validate.Struct(rq); err != nil {
		return err
	}

	if rq.CodeRef == "" {
		rq.CodeRef = rq.Name
	}

	if rq.TestCaseId == "" {
		rq.TestCaseId = rq.CodeRef
	}

	rq.parentUUID = chi.URLParam(r, "itemUuid")

	return nil
}

func (rq *StartTestItemRQ) toItemModel() model.Item {
	return model.Item{
		UUID:        rq.UUID,
		LaunchUUID:  rq.LaunchUUID,
		Name:        rq.Name,
		Description: rq.Description,
		Type:        rq.Type,
		StartTime:   rq.StartTime,
		Attributes:  rq.Attributes.toAttributesModel(),
		Parameters:  rq.Parameters.toParametersModel(),
		CodeRef:     rq.CodeRef,
		TestCaseId:  rq.TestCaseId,
		ParentUUID:  rq.parentUUID,
		IsRetry:     rq.IsRetry,
		RetryOf:     rq.RetryOf,
	}
}

type StartTestItemRS struct {
	UUID string `json:"id"`
}

func (rs *StartTestItemRS) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type FinishTestItemRQ struct {
	LaunchUUID  string           `json:"launchUuid" validate:"required,uuid"`
	EndTime     time.Time        `json:"endTime" validate:"required,datetime"`
	Status      model.ItemStatus `json:"status,omitempty" validate:"omitempty,oneof=PASSED FAILED STOPPED SKIPPED INTERRUPTED CANCELLED INFO WARN"`
	Attributes  Attributes       `json:"attributes,omitempty" validate:"max=256,dive,unique"`
	Description string           `json:"description,omitempty"`
	TestCaseId  string           `json:"testCaseId,omitempty"`
	IsRetry     bool             `json:"retry,omitempty"`
	RetryOf     string           `json:"retryOf,omitempty" validate:"omitempty,uuid"`
	Issue       Issue            `json:"issue,omitempty"`
}

func (rq *FinishTestItemRQ) Bind(_ *http.Request) error {
	if err := validate.Struct(rq); err != nil {
		return err
	}

	if rq.Status == "" {
		rq.Status = model.ItemStatusPassed
	}
	return nil
}

func (rq *FinishTestItemRQ) toItemModel() model.Item {
	return model.Item{
		LaunchUUID:  rq.LaunchUUID,
		EndTime:     &rq.EndTime,
		Status:      rq.Status,
		Attributes:  rq.Attributes.toAttributesModel(),
		Description: rq.Description,
		TestCaseId:  rq.TestCaseId,
	}
}

type FinishTestItemRS struct {
	Message string `json:"message"`
}

func (rs *FinishTestItemRS) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type ItemAttribute struct {
	Key      string `json:"key,omitempty"`
	Value    string `json:"value" validate:"required"`
	IsSystem bool   `json:"system"`
}

type Attributes []ItemAttribute

func (a Attributes) toAttributesModel() model.Attributes {
	attrs := make(model.Attributes, len(a))
	for i, attr := range a {
		attrs[i] = model.Attribute(attr)
	}
	return attrs
}

func fromAttributesModel(model model.Attributes) Attributes {
	attrs := make(Attributes, len(model))
	for i, attr := range model {
		attrs[i] = ItemAttribute(attr)
	}
	return attrs
}

type Parameter struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value,omitempty"`
}

type Parameters []Parameter

func (p Parameters) toParametersModel() model.Parameters {
	params := make(model.Parameters, len(p))
	for i, param := range p {
		params[i] = model.Parameter(param)
	}
	return params
}

func fromParametersModel(model model.Parameters) Parameters {
	params := make(Parameters, len(model))
	for i, param := range model {
		params[i] = Parameter(param)
	}
	return params
}

type Issue struct {
	Type string `json:"issueType" validate:"required"`
}

type TestItemResourceRS struct {
	StartTime time.Time  `json:"startTime"`
	EndTime   *time.Time `json:"endTime,omitempty"`
	TestItemResource
}

type TestItemResourceOldRS struct {
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime,omitempty"`
	TestItemResource
}

func NewTestItemResourceOldRS(item model.Item) *TestItemResourceOldRS {
	var endTime int64
	if item.EndTime != nil {
		endTime = item.EndTime.UnixMilli()
	}

	return &TestItemResourceOldRS{
		StartTime: item.StartTime.UnixMilli(),
		EndTime:   endTime,
		TestItemResource: TestItemResource{
			UUID:        item.UUID,
			Name:        item.Name,
			Type:        item.Type,
			Status:      item.Status,
			CodeRef:     item.CodeRef,
			TestCaseId:  item.TestCaseId,
			Description: item.Description,
			Parameters:  fromParametersModel(item.Parameters),
			Attributes:  fromAttributesModel(item.Attributes),
			Issue:       Issue(item.Issue),
		},
	}
}

func (rs *TestItemResourceOldRS) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type TestItemResource struct {
	UUID        string           `json:"uuid"`
	Name        string           `json:"name"`
	Type        model.ItemType   `json:"type"`
	Status      model.ItemStatus `json:"status"`
	CodeRef     string           `json:"codeRef"`
	TestCaseId  string           `json:"testCaseId"`
	Description string           `json:"description"`
	Parameters  Parameters       `json:"parameters,omitempty"`
	Attributes  Attributes       `json:"attributes,omitempty"`
	Issue       Issue            `json:"issue,omitempty"`
	UndefinedTestItemFields
}

type UndefinedTestItemFields struct {
	ID               int64       `json:"id"`
	LaunchStatus     string      `json:"launchStatus"`
	Parent           int64       `json:"parent"`
	PathNames        interface{} `json:"pathNames"`
	HasStats         bool        `json:"hasStats"`
	HasChildren      bool        `json:"hasChildren"`
	TestCaseHash     int         `json:"testCaseHash"`
	LaunchId         int64       `json:"launchId"`
	UniqueId         string      `json:"uniqueId"`
	PatternTemplates []string    `json:"patternTemplates"`
	Statistics       interface{} `json:"statistics"`
	Path             string      `json:"path"`
}
