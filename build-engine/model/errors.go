package model

import (
	"errors"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

var NoRepositoryConfiguration = errors.New("no repository configuration found")
var NoBuildRun = errors.New("no build run found")