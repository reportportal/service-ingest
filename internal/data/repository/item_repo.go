package repository

import (
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/model"
)

type ItemRepositoryImpl struct {
	buffer buffer.Buffer
}

func NewItemRepository(buffer buffer.Buffer) *ItemRepositoryImpl {
	return &ItemRepositoryImpl{buffer}
}

func (i ItemRepositoryImpl) Get(project string, uuid string) (*model.Item, error) {
	//TODO implement me
	return nil, nil
}

func (i ItemRepositoryImpl) Start(project string, item model.Item) error {
	//TODO implement me
	panic("implement me")
}

func (i ItemRepositoryImpl) Finish(project string, item model.Item) error {
	//TODO implement me
	panic("implement me")
}
