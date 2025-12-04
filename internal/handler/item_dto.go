package handler

import "github.com/reportportal/service-ingest/internal/model"

type ItemAttribute struct {
	Key      string `json:"key,omitempty"`
	Value    string `json:"value" validate:"required"`
	IsSystem bool   `json:"system"`
}

type Attributes []ItemAttribute

func (a Attributes) toAttributesModel() model.Attributes {
	attrs := make([]model.Attribute, 0, len(a))
	for _, a := range a {
		attrs = append(attrs, model.Attribute(a))
	}
	return attrs
}

func fromAttributesModel(modelAttrs model.Attributes) Attributes {
	attrs := make(Attributes, len(modelAttrs))
	for i, attr := range modelAttrs {
		attrs[i] = ItemAttribute(attr)
	}
	return attrs
}
