package handler

import "github.com/reportportal/service-ingest/internal/model"

type ItemAttribute struct {
	Key      string `json:"key,omitempty"`
	Value    string `json:"value"`
	IsSystem bool   `json:"system"`
}

func (a ItemAttribute) toAttributeModel() model.Attribute {
	return model.Attribute{
		Key:      a.Key,
		Value:    a.Value,
		IsSystem: a.IsSystem,
	}
}

func (sl StartLaunchRQ) toAttributesModel() []model.Attribute {
	attrs := make([]model.Attribute, 0, len(sl.Attributes))
	for _, a := range sl.Attributes {
		attrs = append(attrs, a.toAttributeModel())
	}
	return attrs
}
