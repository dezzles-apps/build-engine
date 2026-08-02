package service

import (
	"dezzles-apps/build-engine/model"
	"dezzles-apps/build-engine/repository"
)

func ValidateBuild(event model.BuildEvent) (bool, error) {
	organisation := event.Organisation
	repositoryName := event.Repository
	buildNumber := event.BuildNumber
	ref := event.Ref

	config, err := repository.CreateRepositoryConfiguration(organisation, repositoryName, nil)
	if err != nil {
		return false, err
	}
	if config == nil {
		return false, nil
	}

	build, err := repository.CreateBuildRun(organisation, repositoryName, buildNumber, ref)
	if err != nil {
		return false, err
	}
	if build == nil {
		return false, nil
	}

	return true, nil
}