package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/reportportal/service-ingest/internal/model"
)

type LaunchService struct {
	launchRepo LaunchRepository
}

func NewLaunchService(launchRepo LaunchRepository) *LaunchService {
	return &LaunchService{
		launchRepo: launchRepo,
	}
}

func (s *LaunchService) GetLaunch(project string, uuid string) (model.Launch, error) {
	fmt.Printf("GetLaunch: project=%s uuid=%s\n", project, uuid)

	return model.Launch{}, nil
}

func (s *LaunchService) StartLaunch(project string, launch model.Launch) error {
	launch.Status = model.LaunchStatusInProgress
	launch.UpdatedAt = time.Now().UTC()

	launchJSON, _ := json.MarshalIndent(launch, "", "  ")
	fmt.Printf("Project: %s\nStarting launch:\n%s\n", project, string(launchJSON))

	return nil
}

func (s *LaunchService) FinishLaunch(project string, launchUUID string, launch model.Launch) error {
	launch.UpdatedAt = time.Now().UTC()

	launchJSON, _ := json.MarshalIndent(launch, "", "  ")
	fmt.Printf("Project: %s\nFinishing launch %s:\n%s\n", project, launchUUID, string(launchJSON))

	return nil
}

func (s *LaunchService) UpdateLaunch(project string, launchID int64, launch model.Launch) error {
	launch.UpdatedAt = time.Now().UTC()

	launchJSON, _ := json.MarshalIndent(launch, "", "  ")
	fmt.Printf("Project: %s\nUpdating launch %d:\n%s\n", project, launchID, string(launchJSON))

	return nil
}
