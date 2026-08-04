package model

import "fmt"

type PublishEvent struct {
	Name          string        `json:"name"`
	Version       string        `json:"version"`
	RepositoryUrl string        `json:"repositoryUrl"`
	Config        PublishConfig `json:"config"`
}

type PublishConfig struct {
	Enabled         bool     `json:"enabled"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Url             string   `json:"url"`
	Category        string   `json:"category"`
	ShareRepository bool     `json:"shareRepository"`
	Technologies    []string `json:"technologies"`
}

func (p *PublishEvent) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Version == "" {
		return fmt.Errorf("version is required")
	}
	if p.Config.Title == "" {
		return fmt.Errorf("config.title is required")
	}
	if p.Config.Description == "" {
		return fmt.Errorf("config.description is required")
	}
	if p.Config.Url == "" {
		p.Config.Url = ""
	}
	if p.Config.Category == "" {
		return fmt.Errorf("config.category is required")
	}
	if len(p.Config.Technologies) == 0 {
		return fmt.Errorf("config.technologies must have at least one technology")
	}
	return nil
}
