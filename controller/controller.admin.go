// controller.admin.go

package controller

import (
	"github.com/gin-gonic/gin"
)

func ShowDashboardPage(c *gin.Context) {

	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title":   "Admin Page",
		"content": "Dashboard"}, "admin-dashboard.html")
}
