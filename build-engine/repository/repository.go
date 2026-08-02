package repository

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/go-sql-driver/mysql"
	"dezzles-apps/build-engine/model"
)

var db *sql.DB

func Connect(DatabaseConfig *model.DatabaseConfig) {
	cfg := mysql.NewConfig()

	cfg.User = DatabaseConfig.Username
	cfg.Passwd = DatabaseConfig.Password
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%d", DatabaseConfig.Host, DatabaseConfig.Port)
	cfg.DBName = DatabaseConfig.Database
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
			log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
			log.Fatal(pingErr)
	}
	fmt.Println("Database connected!")
}

func GetBuilds() ([]model.BuildRun, error) {
	rows, err := db.Query("SELECT r.organisation, r.repository, b.build_number, b.start_time FROM builds b JOIN repositories r ON b.repository_id = r.repository_id ORDER BY b.start_time DESC")
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

func GetBuildsByRepository(org string, repository string) ([]model.BuildRun, error) {
	rows, err := db.Query("SELECT r.organisation, r.repository, b.build_number, b.start_time FROM builds b JOIN repositories r ON b.repository_id = r.repository_id WHERE r.organisation = ? AND r.repository = ? ORDER BY b.start_time DESC", org, repository)
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

func GetRepositoryConfiguration(org string, repository string) (*model.RepositoryConfiguration, error) {
	row := db.QueryRow("SELECT organisation, repository, channel FROM repositories WHERE organisation = ? AND repository = ?", org, repository)

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

func GetBuild(org string, repository string, buildNumber int) (*model.DetailedBuildRun, error) {
	row := db.QueryRow("SELECT r.organisation, r.repository, b.build_number, b.start_time, b.discord_thread_id FROM builds b JOIN repositories r ON b.repository_id = r.repository_id WHERE r.organisation = ? AND r.repository = ? AND b.build_number = ?", org, repository, buildNumber)

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

func UpdateDiscordThreadId(org string, repository string, buildNumber int, discordThreadId int64) error {
	// First, get the repositoryId for the given org and repository
	var repositoryId int
	err := db.QueryRow("SELECT repository_id FROM repositories WHERE organisation = ? AND repository = ?", org, repository).Scan(&repositoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.NoRepositoryConfiguration
		}
		return err
	}

	// Now update the discord_thread_id for the given build_number and repositoryId
	result, err := db.Exec("UPDATE builds SET discord_thread_id = ? WHERE build_number = ? AND repository_id = ?", discordThreadId, buildNumber, repositoryId)
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

func CreateRepositoryConfiguration(org string, repository string, channel *string) (*model.RepositoryConfiguration, error) {
	var existing, err = GetRepositoryConfiguration(org, repository)
	if err != nil && err != model.NoRepositoryConfiguration {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	_, err = db.Exec("INSERT INTO repositories (organisation, repository, channel) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE channel = ?", org, repository, channel, channel)
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

func CreateBuildRun(org string, repository string, buildNumber int, ref string) (*model.DetailedBuildRun, error) {
	var existing, err = GetBuild(org, repository, buildNumber)
	if err != nil && err != model.NoBuildRun {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	var repositoryId int
	err = db.QueryRow("SELECT repository_id FROM repositories WHERE organisation = ? AND repository = ?", org, repository).Scan(&repositoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.NoRepositoryConfiguration
		}
		return nil, err
	}

	_, err = db.Exec("INSERT INTO builds (repository_id, build_number, ref) VALUES (?, ?, ?)", repositoryId, buildNumber, ref)
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
