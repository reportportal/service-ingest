package service

import (
	"log/slog"

	"github.com/reportportal/service-ingest/internal/model"
)

type ItemService struct {
	itemRepo ItemRepository
}

func NewItemService(repo ItemRepository) *ItemService {
	return &ItemService{repo}
}

func (s *ItemService) StartItem(project string, item model.Item) error {
	if err := s.itemRepo.Start(project, item); err != nil {
		return err
	}
	slog.Debug("Started item", "project", project, "item", item)
	return nil
}

func (s *ItemService) FinishItem(project string, itemUUID string, item model.Item) error {
	item.UUID = itemUUID
	if err := s.itemRepo.Finish(project, item); err != nil {
		return err
	}
	slog.Debug("Finished item", "project", project, "item", item)
	return nil
}

func (s *ItemService) GetItem(project string, itemUUID string) (model.Item, error) {
	slog.Debug("Getting item", "project", project, "itemUUID", itemUUID)
	return model.Item{
		UUID:       itemUUID,
		LaunchUUID: "example-launch-uuid",
		Name:       "Example Item",
		Type:       model.ItemTypeTest,
		Status:     model.ItemStatusPassed,
	}, nil
}
