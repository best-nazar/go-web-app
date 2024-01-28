// middleware.auth.go

package middleware

import (
	sqladapter "github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
	"github.com/gin-gonic/gin"
)

// This middleware sets whether the user is logged in or not
func SetUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			auth, _ := c.Cookie("auth")
			c.Request.Header.Set("Authorization", auth)
			c.Set("is_logged_in", true) // Used for UI/Menu template (see render() in main.go)

			_, id, errt := helpers.RecoverSessionToken(token)
			c.Set("user_id", int(id))

			if errt != nil {
				c.Set("is_logged_in", false)
			} else {
				cf := c.MustGet("config")
				config := cf.(model.Config)
				if config.UserActivityLogging && c.FullPath() != "" {
					repository.AddUserActivity(c.Request.URL.String(), "path", id)
				}
			}

			c.Next()
		} else {
			c.Set("is_logged_in", false)
			// Set guest user if he's not logged in
			if c.Request.Header.Get("Authorization") == "" {
				c.Request.SetBasicAuth(model.GUEST_ROLE, "")
			}
			c.Next()
		}
	}
}

func CheckCasbinRules() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize an adapter and use it in a Casbin enforcer:
		// the default table name is "casbin_rule".
		// If it doesn't exist, the adapter will create it automatically.
		a, err := sqladapter.NewAdapter(sqladapter.GetDBConnectionInstance())
		if err != nil {
			panic(err)
		}
		// load the casbin model and policy from file "authz_policy.csv", database is also supported.
		casbinEnforcer, err := casbin.NewEnforcer("authz_model.conf", a)
		if err != nil {
			panic(err)
		}

		// You can get acccess to the service from anywhere, by
		// if casbinEnforcer, cExists := c.Get("casbinEnforcer"); cExists {
		// 	casbinEnforcer = casbinEnforcer.(*casbin.Enforcer)
		// }
		c.Set("casbinEnforcer", casbinEnforcer)
		
		authz.NewAuthorizer(casbinEnforcer)(c)
		c.Next()
	}
}