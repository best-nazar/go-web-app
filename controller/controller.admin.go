// controller.admin.go

package controller

import (
	"log"
	"net/http"
	"reflect"

	"github.com/best-nazar/web-app/errorSrc"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
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
		"title":     "Users and Roles",
		"roles":     repository.ListRoles(),
		"activeTab": role,
		"payload":   payload,
	},
		"users-list.html",
		http.StatusOK,
	)
}

func ShowUserRolesPage(c *gin.Context) {
	var casbins interface{}
	var groupRoles interface{}
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
			groupRoles = repository.GetGroupRoles()
		}
	} else {
		// default or not in tabMappings
		tabName = "policy"
		casbins = repository.GetCasbinPolicies()
	}

	Render(c, gin.H{
		"title":   "Users and Roles",
		"page":    "users-roles.html",
		"tab":     tabName,
		"payload": casbins,
		"roles":   groupRoles,
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
			"page":    "users-roles.html",
			"tab":     "role",
			"payload": casbins,
			"error":   errView},
			"admin-dashboard.html",
			http.StatusBadRequest)
		return
	}

	repository.SaveCasbinRole(&role)

	c.Redirect(http.StatusFound, "/admin/uroles?tab=role")
}

func DeleteUserRoles(c *gin.Context) {
	var role = model.CasbinRole{}
	casbins := repository.GetGroupRoles()

	c.ShouldBind(&role)

	recNo, sysErr := repository.DeleteCasbinRole(&role)

	if sysErr != nil {
		errView := errorSrc.MakeErrorView("Delete role", sysErr)

		Render(c, gin.H{
			"title":   "Users and Roles",
			"page":    "users-roles.html",
			"tab":     "role",
			"payload": casbins,
			"error":   errView},
			"admin-dashboard.html",
			http.StatusBadRequest)
		return
	}

	if recNo == 0 {
		errView := errorSrc.MakeErrorViewFrom("Delete role", "ID", http.StatusNotFound)

		Render(c, gin.H{
			"title":   "Users and Roles",
			"page":    "users-roles.html",
			"tab":     "role",
			"payload": casbins,
			"error":   errView},
			"admin-dashboard.html",
			http.StatusNotFound)
		return
	}

	c.Redirect(http.StatusFound, "/admin/uroles?tab=role")
}

func UpdateUserGroups(c *gin.Context) {
	var errView = errorSrc.ErrorView{}
	var jRole = model.CasbinRole{}

	c.ShouldBind(&jRole)

	role, rErr := repository.FindCasbinRolebyName(jRole.Title)
	group, gErr := repository.FindCasbinRoleById(&jRole.ID)

	if rErr != nil {
		errView = errorSrc.MakeErrorViewFrom("Role", "title", http.StatusNotFound)
	} else if gErr != nil {
		errView = errorSrc.MakeErrorViewFrom("Role", "ID", http.StatusNotFound)
	}

	if !reflect.DeepEqual(errView, errorSrc.ErrorView{}) {
		Render(c, gin.H{
			"title":   "Users and Roles",
			"page":    "users-roles.html",
			"tab":     "role",
			"payload": nil,
			"error":   errView},
			"admin-dashboard.html",
			http.StatusNotFound)
		return
	}

	group.V1 = role.Title

	_, nErr := repository.UpdateCusbinRule(group)

	if nErr != nil {
		errView := errorSrc.MakeErrorViewFrom("Role", "ID", http.StatusBadRequest)
		Render(c, gin.H{
			"title":   "Users and Roles",
			"page":    "users-roles.html",
			"tab":     "role",
			"payload": nil,
			"error":   errView},
			"admin-dashboard.html",
			http.StatusBadRequest)
		return
	}

	c.Redirect(http.StatusFound, "/admin/uroles?tab=group")
}

// Action: Remove users from the casbin groups
func RemoveUsersFromGroup(c *gin.Context) {
	var users model.UsersList
	c.ShouldBind(&users)
	repository.DeleteCasbinGroup(&users)

	c.Redirect(http.StatusFound, "/admin/users/list?tab=" + users.Group)
}

// Action: Add user to casbin group
func AddUserToGroup(c *gin.Context) {
	var user model.User
	var userList []string

	e := c.ShouldBind(&user)

	if e != nil {
		log.Println("AddUserToGroup | " + e.Error())
	}

	_, uNum := repository.GetUserByUsername(user.Username)
	roles := repository.ListRoles()
	idx := slices.IndexFunc(*(roles), func(c model.CasbinRole) bool { return c.Title == user.Groups })

	if uNum == 0 {
		errView := errorSrc.MakeErrorViewFrom("Value (" + user.Username + ")", "username", http.StatusNotFound)
		Render(c, gin.H{
			"title":     "Users and Roles",
			"roles":     roles,
			"activeTab": user.Groups,
			"payload":   userList,
			"error":   errView,
		},
			"users-list.html",
			http.StatusNotFound,
		)
		return
	}
	if idx == -1 {
		errView := errorSrc.MakeErrorViewFrom("Value (" + user.Groups + ")", "groups", http.StatusBadRequest)
		Render(c, gin.H{
			"title":     "Users and Roles",
			"roles":     roles,
			"activeTab": user.Groups,
			"payload":   userList,
			"error":   errView,
		},
			"users-list.html",
			http.StatusBadRequest,
		)
		return
	}

	repository.AddCasbinUserRole(user.Username, user.Groups)

	c.Redirect(http.StatusFound, "/admin/users/list?tab=" + user.Groups)
}
