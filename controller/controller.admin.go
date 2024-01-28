// controller.admin.go

package controller

import (
	"net/http"

	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

// use a single instance of Validate, it caches struct info
//var validate *validator.Validate

func ShowDashboardPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title":   "Admin Page",
		"payload": "Dashboard"}, "admin-dashboard.html", http.StatusOK)
}

func ShowUsersListPage(c *gin.Context) {
	var payload []string
	role := c.Query("tab")

	if role == "" {
		role = model.GUEST_ROLE // initiate default tab
	}

	if casbinEnforcer, cExists := c.Get("casbinEnforcer"); cExists {
		casbinEnforcer := casbinEnforcer.(*casbin.Enforcer)
		payload, _ = casbinEnforcer.GetRoleManager().GetUsers(role, "")
	}

	Render(c, gin.H{
		"title":     "Manage Roles",
		"roles":     repository.ListRoles(),
		"activeTab": role,
		"payload":   payload,
	},
		"users-list.html",
		http.StatusOK,
	)
}

// Action: Remove users from the casbin groups
func RemoveUsersFromGroup(c *gin.Context) {
	var users model.UsersList
	e := c.ShouldBind(&users)

	if e == nil {
		repository.DeleteCasbinGroup(&users)
		c.Redirect(http.StatusOK, "/admin/users/list?tab="+users.Group)
	} else {
		c.Redirect(http.StatusFound, "/admin/users/list?tab="+users.Group)
	}
}

// Action: Add user to casbin group
func AddUserToGroup(c *gin.Context) {
	var user model.User
	var userList []string

	e := c.ShouldBind(&user)

	if e != nil {
		c.Error(e)
	}

	_, uNum := repository.GetUserByUsername(user.Username)
	roles := repository.ListRoles()
	idx := slices.IndexFunc(*(roles), func(c model.CasbinRole) bool { return c.Title == user.Groups })

	if uNum == 0 {
		er := errors.New("(" + user.Username + ") not found")
		c.Error(er)

		Render(c, gin.H{
			"title":     "Manage Roles",
			"roles":     &roles,
			"activeTab": user.Groups,
			"payload":   &userList,
			"errors":    c.Errors,
		},
			"users-list.html",
			http.StatusNotFound,
		)
		return
	}

	if idx == -1 {
		er := errors.New("Incorrect value (" + user.Groups + ")")
		c.Error(er)

		Render(c, gin.H{
			"title":     "Manage Roles",
			"roles":     &roles,
			"activeTab": user.Groups,
			"payload":   &userList,
			"errors":    c.Errors,
		},
			"users-list.html",
			http.StatusBadRequest,
		)
		return
	}

	repository.AddCasbinUserRole(user.Username, user.Groups)

	c.Redirect(http.StatusFound, "/admin/users/list?tab="+user.Groups)
}

func ShowCasbinRoutes(c *gin.Context) {
	payload := repository.GetCasbinPolicies()
	roles := repository.ListRoles()

	Render(c, gin.H{
		"title":   "Casbin resources",
		"payload": payload,
		"groups":  roles,
		"actions": model.ACTIONS,
	},
		"casbins-list.html",
		http.StatusOK,
	)
}

func AddCasbinRoute(c *gin.Context) {
	var cr model.CasbinRule

	cr.P_type = model.GROUP_TYPE_P
	cr.V0 = c.PostForm("group")
	cr.V1 = c.PostForm("route")
	cr.V2 = c.PostForm("action")

	print("cr")
}
