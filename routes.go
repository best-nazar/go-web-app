// routes.go

package main

import (
	"github.com/best-nazar/web-app/controller"
)

func initializeRoutes() {
	// Load the APP configuration
	router.Use(setConfiguration())
	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	router.Use(setUserStatus())
	// ACL or RBAC chhecks
	router.Use(checkCasbinRules())

	// Handle the index route
	router.GET("/", controller.ShowIndexPage)

	// Group user related routes together
	userRoutes := router.Group("/u")
	{
		// Handle the GET requests at /u/login
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/login", controller.ShowLoginPage)

		// Handle POST requests at /u/login
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/login", controller.PerformLogin)

		// Handle GET requests at /u/logout
		// Ensure that the user is logged in by using the middleware
		userRoutes.GET("/logout", controller.Logout)

		// Handle the GET requests at /u/register
		// Show the registration page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/register", controller.ShowRegistrationPage)

		// Handle POST requests at /u/register
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/register", controller.Register)
	}

	// Group administrative routes
	adminRoutes := router.Group("admin")
	{
		adminRoutes.GET("/dashboard", controller.ShowDashboardPage)
		adminRoutes.GET(("/uroles"), controller.ShowUserRolesPage)
		adminRoutes.POST(("/uroles"), controller.SaveUserRoles)
		adminRoutes.POST(("/uroles/delete"), controller.DeleteUserRoles)
	}

	// Group article related routes together
	articleRoutes := router.Group("/article")
	{
		// Handle GET requests at /article/view/some_article_id
		articleRoutes.GET("/view/:article_id", controller.GetArticle)

		// Handle the GET requests at /article/create
		// Show the article creation page
		// Ensure that the user is logged in by using the middleware
		articleRoutes.GET("/create", controller.ShowArticleCreationPage)

		// Handle POST requests at /article/create
		// Ensure that the user is logged in by using the middleware
		articleRoutes.POST("/create", controller.CreateArticle)
	}
}