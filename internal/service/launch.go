package service

import (
	"fmt"

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

func (s *LaunchService) StartLaunch(launch model.Launch, project string) (model.Launch, error) {
	fmt.Printf("Project: %+v Starting launch: %+v\n", launch, project)
	return launch, nil
}
