package service

import "github.com/reportportal/service-ingest/internal/model"

type ItemRepository interface {
	Get(project string, uuid string) (*model.Item, error)
	Start(project string, item model.Item) error
	Finish(project string, item model.Item) error
}
