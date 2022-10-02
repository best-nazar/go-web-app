package model

type Config struct {
	DefaultCasbinGroup string `yaml:"default-casbin-group"`
	Tags []string `yaml:"tags"`
}