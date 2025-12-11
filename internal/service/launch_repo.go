package service

import "github.com/reportportal/service-ingest/internal/model"

type LaunchRepository interface {
	Get(project string, uuid string) (*model.Launch, error)
	Create(project string, launch model.Launch) error
	Update(project string, launch model.Launch) error
}
