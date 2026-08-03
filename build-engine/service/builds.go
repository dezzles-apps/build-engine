package service

import (
	"dezzles-apps/build-engine/model"
	"dezzles-apps/build-engine/repository"
)

type BuildService struct {
	repository *repository.EventsRepository
}

func (s *BuildService) Initialise(repository *repository.EventsRepository) error {
	s.repository = repository
	return nil
}

func (s *BuildService) ValidateBuild(event model.BuildEvent) (bool, error) {
	organisation := event.Organisation
	repositoryName := event.Repository
	buildNumber := event.BuildNumber
	ref := event.Ref

	config, err := s.repository.CreateRepositoryConfiguration(organisation, repositoryName, nil)
	if err != nil {
		return false, err
	}
	if config == nil {
		return false, nil
	}

	build, err := s.repository.CreateBuildRun(organisation, repositoryName, buildNumber, ref)
	if err != nil {
		return false, err
	}
	if build == nil {
		return false, nil
	}

	return true, nil
}
