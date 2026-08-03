package repository

import (
	"database/sql"
	"dezzles-apps/build-engine/model"
	"fmt"
)

type EventsRepository struct {
	database *Database
}

func (r *EventsRepository) Initialise(database *Database) error {
	r.database = database
	return nil
}

func (r *EventsRepository) GetBuilds() ([]model.BuildRun, error) {
	rows, err := r.database.getDB().Query("SELECT r.organisation, r.repository, b.build_number, b.start_time FROM builds b JOIN repositories r ON b.repository_id = r.repository_id ORDER BY b.start_time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var builds []model.BuildRun
	for rows.Next() {
		var build model.BuildRun
		if err := rows.Scan(&build.Organisation, &build.Repository, &build.BuildNumber, &build.StartTime); err != nil {
			return nil, err
		}
		builds = append(builds, build)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return builds, nil
}

func (r *EventsRepository) GetBuildsByRepository(org string, repository string) ([]model.BuildRun, error) {
	rows, err := r.database.getDB().Query("SELECT r.organisation, r.repository, b.build_number, b.start_time FROM builds b JOIN repositories r ON b.repository_id = r.repository_id WHERE r.organisation = ? AND r.repository = ? ORDER BY b.start_time DESC", org, repository)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var builds []model.BuildRun
	for rows.Next() {
		var build model.BuildRun
		if err := rows.Scan(&build.Organisation, &build.Repository, &build.BuildNumber, &build.StartTime); err != nil {
			return nil, err
		}
		builds = append(builds, build)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return builds, nil
}

func (r *EventsRepository) GetRepositoryConfiguration(org string, repository string) (*model.RepositoryConfiguration, error) {
	row := r.database.getDB().QueryRow("SELECT organisation, repository, channel FROM repositories WHERE organisation = ? AND repository = ?", org, repository)

	var config model.RepositoryConfiguration
	var nullChannel sql.NullString
	if err := row.Scan(&config.Organisation, &config.Repository, &nullChannel); err != nil {
		if err == sql.ErrNoRows {
			return nil, model.NoRepositoryConfiguration
		}
		return nil, err
	}
	if nullChannel.Valid {
		config.Channel = &nullChannel.String
	} else {
		defaultChannel := "dezzles-apps"
		config.Channel = &defaultChannel
	}
	return &config, nil
}

func (r *EventsRepository) GetBuild(org string, repository string, buildNumber int) (*model.DetailedBuildRun, error) {
	row := r.database.getDB().QueryRow("SELECT r.organisation, r.repository, b.build_number, b.start_time, b.discord_thread_id FROM builds b JOIN repositories r ON b.repository_id = r.repository_id WHERE r.organisation = ? AND r.repository = ? AND b.build_number = ?", org, repository, buildNumber)

	var build model.DetailedBuildRun
	var nullDiscordThreadId sql.NullInt64
	if err := row.Scan(&build.Organisation, &build.Repository, &build.BuildNumber, &build.StartTime, &nullDiscordThreadId); err != nil {
		if err == sql.ErrNoRows {
			return nil, model.NoBuildRun
		}
		return nil, err
	}
	if nullDiscordThreadId.Valid {
		build.DiscordThreadId = &nullDiscordThreadId.Int64
	} else {
		build.DiscordThreadId = nil
	}
	return &build, nil
}

func (r *EventsRepository) UpdateDiscordThreadId(org string, repository string, buildNumber int, discordThreadId int64) error {
	// First, get the repositoryId for the given org and repository
	var repositoryId int
	err := r.database.getDB().QueryRow("SELECT repository_id FROM repositories WHERE organisation = ? AND repository = ?", org, repository).Scan(&repositoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.NoRepositoryConfiguration
		}
		return err
	}

	// Now update the discord_thread_id for the given build_number and repositoryId
	result, err := r.database.getDB().Exec("UPDATE builds SET discord_thread_id = ? WHERE build_number = ? AND repository_id = ?", discordThreadId, buildNumber, repositoryId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no build found for %s/%s with build number %d", org, repository, buildNumber)
	}

	return nil
}

func (r *EventsRepository) CreateRepositoryConfiguration(org string, repository string, channel *string) (*model.RepositoryConfiguration, error) {
	var existing, err = r.GetRepositoryConfiguration(org, repository)
	if err != nil && err != model.NoRepositoryConfiguration {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	_, err = r.database.getDB().Exec("INSERT INTO repositories (organisation, repository, channel) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE channel = ?", org, repository, channel, channel)
	if err != nil {
		return nil, err
	}

	config := &model.RepositoryConfiguration{
		Organisation: org,
		Repository:   repository,
		Channel:      channel,
	}
	return config, nil
}

func (r *EventsRepository) CreateBuildRun(org string, repository string, buildNumber int, ref string) (*model.DetailedBuildRun, error) {
	var existing, err = r.GetBuild(org, repository, buildNumber)
	if err != nil && err != model.NoBuildRun {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	var repositoryId int
	err = r.database.getDB().QueryRow("SELECT repository_id FROM repositories WHERE organisation = ? AND repository = ?", org, repository).Scan(&repositoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.NoRepositoryConfiguration
		}
		return nil, err
	}

	_, err = r.database.getDB().Exec("INSERT INTO builds (repository_id, build_number, ref) VALUES (?, ?, ?)", repositoryId, buildNumber, ref)
	if err != nil {
		return nil, err
	}

	build := &model.DetailedBuildRun{
		Organisation: org,
		Repository:   repository,
		BuildNumber:  buildNumber,
		StartTime:    "", // You might want to fetch the actual start time from the database if needed
	}
	return build, nil
}
