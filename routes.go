// routes.go

package main

import (
	"github.com/best-nazar/web-app/controller"
	"github.com/best-nazar/web-app/middleware"
)

func initializeRoutes() {
	// Load the APP configuration
	router.Use(middleware.SetConfiguration())
	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	router.Use(middleware.SetUserStatus())
	// ACL or RBAC checks
	router.Use(middleware.CheckCasbinRules())

	//the index page
	router.GET("/", controller.ShowIndexPage)

	guestRoutes := router.Group("/u")
	{
		guestRoutes.GET("/", controller.ShowUserHomePage)
		// Show the login page
		guestRoutes.GET("/login", controller.ShowLoginPage)
		// Ensure that the user is not logged in by using the middleware
		guestRoutes.POST("/login", controller.PerformLogin)
		// Perform logout
		guestRoutes.GET("/logout", controller.Logout)
		// Show the registration page
		guestRoutes.GET("/register", controller.ShowRegistrationPage)
		// Redirect user after registration succeeded.
		guestRoutes.GET("/register/success", controller.ShowRegistrationSuccess)
		// Submit registration data
		guestRoutes.POST("/register", controller.Register)
		// INFO page. User locked
		guestRoutes.GET("/locked", controller.UserLocked)
	}

	// Action for user within a member area access
	memberRoutes := router.Group("member")
	{
		memberRoutes.GET("/avatar/:id", controller.UploadImage)
		memberRoutes.POST("/avatar/:id", controller.UploadImage)
	}

	// Group administrative routes
	adminRoutes := router.Group("admin")
	{
		adminRoutes.GET("/dashboard", controller.ShowDashboardPage)
		
		adminRoutes.GET("/groups/list", controller.ShowGroupsListPage)
		adminRoutes.POST("/groups/remove", controller.RemoveUserFromGroup)
		adminRoutes.POST("/groups/add", controller.AddUserToGroup)

		adminRoutes.GET("/casbins/list", controller.ShowCasbinRoutes)
		adminRoutes.POST("/casbins/add", controller.AddCasbinRoute)
		adminRoutes.POST("/casbins/remove", controller.RemoveCasbinRoute)

		adminRoutes.GET("/users/list", controller.UsersList)
		adminRoutes.GET("/user/details/:id", controller.UserDetails)
		adminRoutes.POST("/user/update", controller.UserUpdate)
		adminRoutes.POST("/user/update/status", controller.UserActivateDeactivate)
	}
}