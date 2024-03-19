// controller.admin.go

package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/model/dto"
	"github.com/best-nazar/web-app/repository"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// Default page for logged user
func ShowDashboardPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title":   "Admin Page",
		"payload": "Dashboard"}, "admin-dashboard.html", http.StatusOK)
}

// Manage Users and Groups page
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
	var ug dto.UserGroup

	e := c.ShouldBind(&ug)

	if e != nil {
		c.Error(e)
	}

	validateRoles(c, ug.Group)

	_, uNum := repository.GetUserByUsername(ug.Username)

	if uNum == 0 {
		c.Error(fmt.Errorf("username '%s' not found", ug.Username))
	}

	if _, count := repository.FindCasbinGroupByNameAndRole(ug.Username, ug.Group); count != 0 {
		c.Error(fmt.Errorf("username '%s' in Group '%s' already exist", ug.Username, ug.Group))
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
		c.Redirect(http.StatusFound, "/admin/groups/list?tab=" + ug.Group)
	}
}

// Manage URL resources and access groups page
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

// Manage URL resources page
// Add route with access permissions
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
			c.Error(fmt.Errorf("rule '" + cr.V1 + " " + cr.V0 + " " + cr.V2 + "' already exists"))
		}
	} else {
		c.Error(fmt.Errorf("route must be valid URL string"))
	}

	validateRoles(c, cr.V0)

	if idx := slices.Index(model.ACTIONS, cr.V2); idx == -1 {
		c.Error(fmt.Errorf("action '%s' is not allowed. Allowed values: %s", cr.V2, strings.Join(model.ACTIONS, ", ")))
	}

	if len(c.Errors) > 0 {
		ShowCasbinRoutes(c)
	} else {
		repository.AddCasbinRole(&cr)
		c.Redirect(http.StatusFound, "/admin/casbins/list")
	}
}

// Manage URL resources page
// Remove route from Casbin management 
func RemoveCasbinRoute(c *gin.Context) {
	c.Request.ParseForm()
	for key, values := range c.Request.PostForm {
		if key != "ID" {
			c.Error(fmt.Errorf("missing ID"))
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

// Show the full list of registered Users page
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

// User details page
func UserDetails(c *gin.Context) {
	id := c.Param("id")

	user, nRows := repository.FindUserById(id)

	if nRows > 0 {
		casbinEnforcer := c.MustGet("casbinEnforcer").(*casbin.Enforcer)
		groups, _ := casbinEnforcer.GetRolesForUser(user.Username, "")

		avatar_path := ""
		if avatar := user.Avatar(); avatar == nil {
			avatar_path = model.DEFAULT_AVATAR
		} else {
			avatar_path = avatar.Path
		}

		Render(c, gin.H{
			"title":       "User Details",
			"description": user.Name,
			"payload":     user,
			"avatar_path": avatar_path,
			"groups": 	   strings.Join(groups, ", "),
			"errors":      helpers.Errors(c),
		},
			"user-details.html",
			http.StatusOK,
		)
	} else {
		c.Error(fmt.Errorf("user ID '%s' not found", id))

		UsersList(c)
	}
}

// Action: Update User's data
func UserUpdate(c *gin.Context) {
	var user *dto.UpdateUserDto

	err := c.ShouldBind(&user)

	if err == nil {
		repository.UpdateUser(user)
		c.Redirect(http.StatusFound, "/admin/user/details/" + user.ID)
	} else {
		c.Error(err)
		UsersList(c)
	}
}

// Action: make User inactive (locked) or vise-versa
func UserActivateDeactivate(c *gin.Context) {
	user := &dto.ActivateUserDto{}
	err := c.ShouldBind(user)

	if (err != nil) {
		c.Error(err)
	}

	findUser, rNum := repository.FindUserById(user.ID)

	if (rNum == 0) {
		c.Error(fmt.Errorf("user ID '%s' not found", user.ID))
	} else {
		if findUser.Active == 0 {
			findUser.Active = 1
			repository.ActivateUser(findUser)
		} else {  
			findUser.Active = 0
			currentTime := time.Now()
			findUser.SuspendedAt = &currentTime

			repository.DeactivateUser(findUser)
		}
	}
	
	c.Redirect(http.StatusFound, "/admin/user/details/" + fmt.Sprintf("%v", findUser.ID))
}

func validateRoles(c *gin.Context, group string) {
	roles := repository.ListRoles()
	idx := slices.IndexFunc(*(roles), func(c model.CasbinRole) bool { return c.Title == group })

	if idx == -1 {
		c.Error(fmt.Errorf("group '%s' not found", group))
	}
}
