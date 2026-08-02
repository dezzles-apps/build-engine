package model

type BuildEvent struct {
	Repository string `json:"repository"`
	Organisation string `json:"organisation"`
	BuildNumber int `json:"buildNumber"`
	Message string `json:"message"`
	Component string `json:"component"`
	Ref string `json:"ref"`
}

func (e *BuildEvent) Validate() error {
	if e.Organisation == "" {
		return &ValidationError{Field: "organisation", Message: "Missing required field: organisation"}
	}
	if e.Repository == "" {
		return &ValidationError{Field: "repository", Message: "Missing required field: repository"}
	}
	if e.BuildNumber == 0 {
		return &ValidationError{Field: "buildNumber", Message: "Missing required field: buildNumber"}
	}
	if e.Message == "" {
		return &ValidationError{Field: "message", Message: "Missing required field: message"}
	}
	if e.Ref == "" {
		return &ValidationError{Field: "ref", Message: "Missing required field: ref"}
	}
	if e.Component == "" {
		e.Component = "default"
	}
	return nil
}