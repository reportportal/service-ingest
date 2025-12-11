package service

import "github.com/reportportal/service-ingest/internal/model"

type ItemRepository interface {
	Get(project string, uuid string) (*model.Item, error)
	Create(project string, item model.Item) error
	Update(project string, item model.Item) error
}
