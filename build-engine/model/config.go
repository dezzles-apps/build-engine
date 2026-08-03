package model

import (
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database     DatabaseConfig     `yaml:"database"`
	EventService EventServiceConfig `yaml:"event-service"`
	Server       ServerConfig       `yaml:"server"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type EventServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// resolveEnv returns the environment variable value if s starts with $, otherwise returns s as-is
func resolveEnv(s string) string {
	if strings.HasPrefix(s, "$") {
		envVar := strings.TrimPrefix(s, "$")
		return os.Getenv(envVar)
	}
	return s
}

// resolveEnvVarsInStruct recursively walks through a struct and resolves environment variables in all string fields
func resolveEnvVarsInStruct(v interface{}) {
	val := reflect.ValueOf(v)

	// Handle pointer types
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Only process structs
	if val.Kind() != reflect.Struct {
		return
	}

	// Iterate through all fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// If it's a pointer, dereference it
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		switch field.Kind() {
		case reflect.String:
			// Resolve environment variables in string fields
			resolved := resolveEnv(field.String())
			field.SetString(resolved)
		case reflect.Struct:
			// Recursively process nested structs
			resolveEnvVarsInStruct(field.Addr().Interface())
		}
	}
}

func LoadConfig() (*Config, error) {
	var config Config
	env := os.Getenv("ENVIRONMENT")
	configFile := "config/config." + env + ".yaml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Automatically resolve all environment variables in string fields
	resolveEnvVarsInStruct(&config)

	return &config, nil
}
