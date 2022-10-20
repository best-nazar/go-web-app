package model

type Config struct {
	DefaultCasbinGroup string `yaml:"default-casbin-group"`
	UserActivityLogging bool `yaml:"user-activity-logging"`
	Tags []string `yaml:"tags"`
}