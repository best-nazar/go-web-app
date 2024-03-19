package model

type Config struct {
	DefaultCasbinGroup 		string 					`yaml:"default_casbin_group"`
	UserActivityLogging 	bool 					`yaml:"user_activity_logging"`
	UsernameRestrictedWords string 					`yaml:"username_restricted_words"`
	Tags 					[]string 				`yaml:"tags"`
	ContactSupportEmail		string					`yaml:"contact_support_email"`
	ImageConfig 			map[string]interface{} 	`yaml:"image_config"`
}