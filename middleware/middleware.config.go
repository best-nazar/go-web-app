package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"github.com/best-nazar/web-app/model"
)

// Sets configuration to the gin.c
func SetConfiguration() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := model.Config{}
		data := loadYamlFile()
		err := yaml.Unmarshal([]byte(data), &conf)

		if err != nil {
			panic(err)
		}
		
		c.Set("config", conf)
		c.Next()
	}
}
// Load the configuration from YAML
func loadYamlFile() []byte {
	path := getPath()
	data, err := os.ReadFile(path + "/config/config.yaml")

	if err != nil {
		panic(err)
	}

	return data
}

// Gets main folder path
func getPath () string {
	ex, err := os.Getwd()
    if err != nil {
        panic(err)
    }

	return ex
}