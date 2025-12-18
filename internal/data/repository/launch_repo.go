package repository

import (
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type LaunchRepositoryImpl struct {
	buffer buffer.Buffer
}

func NewLaunchRepo(buffer buffer.Buffer) *LaunchRepositoryImpl {
	return &LaunchRepositoryImpl{buffer}
}

func (l LaunchRepositoryImpl) Get(project string, uuid string) (*model.Launch, error) {
	//TODO implement me
	return nil, nil
}

func (l LaunchRepositoryImpl) Start(project string, launch model.Launch) error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (l LaunchRepositoryImpl) Update(project string, launch model.Launch) error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (l LaunchRepositoryImpl) Finish(project string, launch model.Launch) error {
	//TODO implement me
	panic("implement me")
	return nil
}
