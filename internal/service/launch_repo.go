package service

import "github.com/reportportal/service-ingest/internal/model"

type LaunchRepository interface {
	Create(launch model.Launch) error
	GetByUUID(uuid string) (*model.Launch, error)
	Update(launch model.Launch) error
}
