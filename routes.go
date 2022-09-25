// routes.go

package main

import (
	"github.com/best-nazar/web-app/controller"
	sqladapter "github.com/best-nazar/web-app/db"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
)

func initializeRoutes() {
	// Initialize an adapter and use it in a Casbin enforcer:
	// the default table name is "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	a, err := sqladapter.NewAdapter(sqladapter.GetDBConnectionInstance())
	if err != nil {
		panic(err)
	}
	// load the casbin model and policy from file "authz_policy.csv", database is also supported.
	e, err := casbin.NewEnforcer("authz_model.conf", a)
	if err != nil {
		panic(err)
	}

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	router.Use(setUserStatus())
	router.Use(authz.NewAuthorizer(e))

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
