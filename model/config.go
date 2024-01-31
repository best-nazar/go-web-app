package model

type Config struct {
	DefaultCasbinGroup string `yaml:"default-casbin-group"`
	UserActivityLogging bool `yaml:"user-activity-logging"`
	UsernameRestrictedWords string `yaml:"username-restricted-words"`
	Tags []string `yaml:"tags"`
}