// main.go

package main

import (
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
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/*/*.html")

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run()
}
