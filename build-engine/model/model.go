package model

type BuildRun struct {
	Organisation string `json:"organisation"`
	Repository string `json:"repository"`
	BuildNumber int `json:"buildNumber"`
	StartTime string `json:"startTime"`
}

type RepositoryConfiguration struct {
	Organisation string `json:"organisation"`
	Repository string `json:"repository"`
	Channel *string `json:"channel"`
}

type DetailedBuildRun struct {
	Organisation string `json:"organisation"`
	Repository string `json:"repository"`
	BuildNumber int `json:"buildNumber"`
	StartTime string `json:"startTime"`
	DiscordThreadId *int64 `json:"discordThreadId"`
}