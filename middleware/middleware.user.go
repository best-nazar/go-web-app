package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/best-nazar/web-app/model"
	"github.com/gin-gonic/gin"
	"github.com/casbin/casbin/v2"
)

var urlToExclude = []string {
	"/u/locked",
	"/u/logout",
}

func IsUserActive() gin.HandlerFunc {
	return func(c *gin.Context) {
		usr := c.MustGet("user")
		path := c.Request.URL.Path
		// default
		c.Set("user_groups", []string{})

		if usr != nil && !slices.Contains(urlToExclude, path) {
			user := usr.(*model.User)

			casbinEnforcer := c.MustGet("casbinEnforcer").(*casbin.Enforcer)
			groups, _ := casbinEnforcer.GetRolesForUser(user.Username, "")
			c.Set("user_groups", groups)

			if user.Active == 0 {
				c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/u/locked?id=%v", user.ID))
				c.Abort()
			}
		}

		c.Next()
	}
}