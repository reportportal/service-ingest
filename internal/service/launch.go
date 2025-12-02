package service

import "github.com/reportportal/service-ingest/internal/model"

type LaunchService struct {
	launchRepo LaunchRepository
}

func NewLaunchService(launchRepo LaunchRepository) *LaunchService {
	return &LaunchService{
		launchRepo: launchRepo,
	}
}

func (s *LaunchService) StartLaunch(launch model.Launch) (model.Launch, error) {
	println("Starting launch:", launch)
	return launch, nil
}
