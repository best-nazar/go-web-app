// controller.admin.go

package controller

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

func ShowDashboardPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title":   "Admin Page",
		"payload": "Dashboard"}, "admin-dashboard.html", http.StatusOK)
}

func ShowGroupsListPage(c *gin.Context) {
	var payload []string
	role := c.Query("tab")

	if role == "" {
		role = model.GUEST_ROLE // initiate default tab
	}

	casbinEnforcer := c.MustGet("casbinEnforcer").(*casbin.Enforcer)
	payload, _ = casbinEnforcer.GetRoleManager().GetUsers(role, "")

	Render(c, gin.H{
		"title":     "Manage Groups",
		"roles":     repository.ListRoles(),
		"activeTab": role,
		"payload":   payload,
	},
		"groups-list.html",
		http.StatusOK,
	)
}

// Action: Remove users from the casbin groups
func RemoveUserFromGroup(c *gin.Context) {
	var users model.UsersList
	e := c.ShouldBind(&users)

	if e == nil {
		repository.DeleteCasbinGroup(&users)
		c.Redirect(http.StatusFound, "/admin/groups/list?tab="+users.Group)
	} else {
		c.Redirect(http.StatusNotModified, "/admin/groups/list?tab="+users.Group)
	}
}

// Action: Add user to casbin group
func AddUserToGroup(c *gin.Context) {
	var ug model.UserGroup

	e := c.ShouldBind(&ug)

	if e != nil {
		c.Error(e)
	}

	validateRoles(c, ug.Group)

	_, uNum := repository.GetUserByUsername(ug.Username)

	if uNum == 0 {
		er := errors.New("username|" + ug.Username + " not found")
		c.Error(er)
	}

	if _, count := repository.FindCasbinGroupByNameAndRole(ug.Username, ug.Group); count != 0 {
		c.Error(errors.New("username|" + ug.Username + " in Group " + ug.Group + " already exist"))
	}

	if len(c.Errors) > 0 {
		casbinEnforcer := c.MustGet("casbinEnforcer").(*casbin.Enforcer)
		payload, _ := casbinEnforcer.GetRoleManager().GetUsers(ug.Group, "")
		roles := repository.ListRoles()

		Render(c, gin.H{
			"title":     "Manage Roles",
			"roles":     roles,
			"activeTab": ug.Group,
			"payload":   payload,
			"errors":    helpers.Errors(c),
		},
			"groups-list.html",
			http.StatusConflict,
		)
	} else {
		repository.AddCasbinUserRole(ug.Username, ug.Group)
		c.Redirect(http.StatusFound, "/admin/groups/list?tab="+ug.Group)
	}
}

func ShowCasbinRoutes(c *gin.Context) {
	payload := repository.GetCasbinPolicies()
	roles := repository.ListRoles()

	Render(c, gin.H{
		"title":   "Manage URL resources",
		"payload": payload,
		"groups":  roles,
		"actions": model.ACTIONS,
		"errors":  helpers.Errors(c),
	},
		"casbins-list.html",
		http.StatusOK,
	)
}

func AddCasbinRoute(c *gin.Context) {
	var cr model.CasbinRuleP

	e := c.ShouldBind(&cr)

	if e != nil {
		c.Error(e)
	}

	u, err := url.ParseRequestURI(cr.V1)

	if err == nil {
		cr.V1 = u.Path
		_, ferr := repository.FindCasbinUrlGroup(&cr)

		if ferr == nil {
			c.Error(errors.New("rules|" + cr.V1 + " " + cr.V0 + " " + cr.V2 + " exists"))
		}
	} else {
		c.Error(errors.New("route|must be valid URL string"))
	}

	validateRoles(c, cr.V0)

	if idx := slices.Index(model.ACTIONS, cr.V2); idx == -1 {
		c.Error(errors.New("action|" + cr.V2 + " is not allowed. Allowed values:" + strings.Join(model.ACTIONS, ", ")))
	}

	if len(c.Errors) > 0 {
		ShowCasbinRoutes(c)
	} else {
		repository.AddCasbinRole(&cr)
		c.Redirect(http.StatusFound, "/admin/casbins/list")
	}
}

func RemoveCasbinRoute(c *gin.Context) {
	c.Request.ParseForm()
	for key, values := range c.Request.PostForm {
		if key != "ID" {
			c.Error(errors.New("Missing ID"))
			c.AbortWithError(http.StatusBadRequest, c.Errors.Last())
			return
		}

		err := repository.RemoveCasbinRole(values)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	c.Redirect(http.StatusFound, "/admin/casbins/list")
}

func UsersList(c *gin.Context) {
	statusCode := http.StatusOK

	if len(c.Errors)>0 {
		statusCode = http.StatusBadRequest
	}

	Render(c, gin.H{
		"title":   "Registered Users",
		"payload": repository.GetUsers(),
		"errors":  helpers.Errors(c),
	},
		"users-list.html",
		statusCode,
	)
}

func UserDetails(c *gin.Context) {
	id := c.Param("id")

	user, nRows := repository.FindUserById(id)

	if nRows > 0 {
		Render(c, gin.H{
			"title":       "User Details",
			"description": user.Name,
			"payload":     user,
			"errors":      helpers.Errors(c),
		},
			"user-details.html",
			http.StatusOK,
		)
	} else {
		c.Error(errors.New("user|ID:" + id + " not found"))

		UsersList(c)
	}
}

func UserUpdate(c *gin.Context) {
	var user *model.UpdateUser

	err := c.ShouldBind(&user)

	if err == nil {
		repository.UpdateUser(user)
		UsersList(c)
	} else {
		c.Error(err)
		UsersList(c)
	}
}

func validateRoles(c *gin.Context, group string) {
	roles := repository.ListRoles()
	idx := slices.IndexFunc(*(roles), func(c model.CasbinRole) bool { return c.Title == group })

	if idx == -1 {
		er := errors.New("group|" + group + " not found")
		c.Error(er)
	}
}
