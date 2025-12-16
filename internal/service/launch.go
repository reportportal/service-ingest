package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/reportportal/service-ingest/internal/model"
)

type LaunchService struct {
	launchRepo LaunchRepository
}

func NewLaunchService(repo LaunchRepository) *LaunchService {
	return &LaunchService{repo}
}

func (s *LaunchService) GetLaunch(project string, uuid string) (*model.Launch, error) {
	launch, err := s.launchRepo.Get(project, uuid)
	if err != nil {
		return nil, err
	}

	return launch, nil
}

func (s *LaunchService) StartLaunch(project string, launch model.Launch) (string, error) {
	launch.UpdatedAt = time.Now().UTC()
	launch.Status = model.LaunchStatusInProgress

	if launch.UUID == "" {
		launch.UUID = uuid.New().String()
	}

	if err := s.launchRepo.Create(project, launch); err != nil {
		return "", err
	}

	return launch.UUID, nil
}

func (s *LaunchService) FinishLaunch(project string, launchUUID string, launch model.Launch) error {
	launch.UpdatedAt = time.Now().UTC()
	launch.UUID = launchUUID

	if err := s.launchRepo.Update(project, launch); err != nil {
		return err
	}

	return nil
}

func (s *LaunchService) UpdateLaunch(project string, launchID int64, launch model.Launch) error {
	launch.UpdatedAt = time.Now().UTC()

	if err := s.launchRepo.Update(project, launch); err != nil {
		return err
	}

	return nil
}
