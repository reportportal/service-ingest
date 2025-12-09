package service

import (
	"fmt"

	"github.com/reportportal/service-ingest/internal/model"
)

type ItemService struct{}

func NewItemService() *ItemService {
	return &ItemService{}
}

func (s *ItemService) StartItem(project string, item model.Item) error {
	fmt.Printf("Project: %s\nSaving item:\n%+v\n", project, item)
	return nil
}

func (s *ItemService) FinishItem(project string, itemUUID string, item model.Item) error {
	fmt.Printf("Project: %s\nFinishing item %s:\n%+v\n", project, itemUUID, item)
	return nil
}

func (s *ItemService) GetItem(project string, itemUUID string) (model.Item, error) {
	fmt.Printf("Project: %s\nGetting item %s\n", project, itemUUID)
	return model.Item{
		UUID:       itemUUID,
		LaunchUUID: "example-launch-uuid",
		Name:       "Example Item",
		Type:       model.ItemTypeTest,
		Status:     model.ItemStatusPassed,
	}, nil
}
