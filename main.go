// main.go

package main

import (
	"errors"
	"html/template"

	"github.com/best-nazar/web-app/helpers"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.DebugMode)

	// Set the router as the default one provided by Gin
	router = gin.Default()
	// Loading static assets like JS & CSS
	router.Static("/assets", "./assets")
	router.Static("/images", "./images")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.SetFuncMap(template.FuncMap{
        "dict": templateDict,
		"formatDate": helpers.TimestampToSting,
    })

	router.LoadHTMLGlob("templates/*/*.html")

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run()
}

// Custom function to pass data inside Templates.
// Usage: {{ template "common-confirmation.html" dict "post_url" "?go/there" }}
// Usege: {{template "userlist.html" dict "Users" .MostPopular "Current" .CurrentUser}}
func templateDict (values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i+=2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}