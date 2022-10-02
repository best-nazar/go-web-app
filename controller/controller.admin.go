// controller.admin.go

package controller

import (
	"net/http"
	"strings"

	"github.com/best-nazar/web-app/errorSrc"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/gin-gonic/gin"
)

// use a single instance of Validate, it caches struct info
//var validate *validator.Validate

func ShowDashboardPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title":   "Admin Page",
		"payload": "Dashboard"}, "admin-dashboard.html", http.StatusOK)
}

func ShowUserRolesPage(c *gin.Context) {
	var casbins interface{}
	tabMappings := map[string]string{"policy": "p", "role": "r", "group": "g"}
	tabName := c.Query("tab")

	if val, ok := tabMappings[tabName]; ok {
		switch val {
		case "p": 
			casbins = repository.GetCasbinPolicies()
		case "r":
			casbins = repository.GetGroupRoles()
		case "g":
			casbins = repository.GetCasbinRoles()
		}
	} else {
		// default or not in tabMappings
		tabName = "policy"
		casbins = repository.GetCasbinPolicies()
	}

	Render(c, gin.H{
		"title":   "Users and Roles",
		"page": "users-roles.html",
		"tab": tabName,
		"payload": casbins,
		}, "admin-dashboard.html", http.StatusOK)
}

func SaveUserRoles(c *gin.Context) {
	var role model.CasbinRole
	err := c.ShouldBind(&role)

	if err != nil {
		casbins := repository.GetGroupRoles()
		errView := errorSrc.MakeErrorView("Add role", err)

		Render(c, gin.H{
			"title":   "Users and Roles",
			"page": "users-roles.html",
			"tab": "role",
			"payload": casbins,
			"error": errView},
			"admin-dashboard.html",
			http.StatusBadRequest)
		return
	}
	// If it's a UI Form request, we need convert array of string to comma separated strings
	postFormInheritance := c.PostFormArray("inheritedFrom")
	if len(postFormInheritance) >0 {
		role.InheritedFrom = strings.Join(postFormInheritance, ",")
	}

	repository.SaveCasbinRole(&role)

	c.Redirect(http.StatusFound, "/admin/uroles?tab=role")
}