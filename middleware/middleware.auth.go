// middleware.auth.go

package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	sqladapter "github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
	"github.com/gin-gonic/gin"
)

var urlToExclude = []string {
	"/u/locked",
	"/u/logout",
}

// This middleware sets whether the user is logged in or not
func SetUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *model.User
		var err error		

		username, password, isAuth := c.Request.BasicAuth() // if request comes from API
		token, _ := c.Cookie("token") //if request comes from UI


		if isAuth {
			user, err = getUserFromAuth(username)

			if err != nil || !user.IsPasswordValid(password) {
				c.AbortWithStatus(http.StatusForbidden)
			}
		} else {
			ipAddr := c.ClientIP()

			user, err = getUserFromToken(token, ipAddr)
		}
		
		if err == nil {
			c.Set("is_logged_in", true) // Used for UI/Menu template (see render.go)
			c.Set("user", user) //UpdateUser struct has no password property for security reason

			isUserActive(c, user)
			//casbin RBAC uses BasicAuth() to get user for rule validation
			c.Request.SetBasicAuth(user.Username, "") 
			config := c.MustGet("config").(model.Config)

			if config.UserActivityLogging && c.FullPath() != "" {
				repository.AddUserActivity(c.Request.URL.String(), "path", user.ID.String())
			}

			c.Next()
		} else {
			// user was not found in token. let's search in BasicAuth
			c.Set("is_logged_in", false)
			c.Set("user", nil)
			// Set guest user if he's not logged in for casbin auth (we use 'guest' for public url)
			c.Request.SetBasicAuth(model.GUEST_ROLE, "")

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
		setUserGroups(c)
		c.Next()
	}
}

func getUserFromToken(token string, ip string) (*model.User, error) {
	_, id, errt := helpers.RecoverSessionToken(token, ip)

	if errt != nil {
		return nil, errt
	}
	
	user, num := repository.FindUserById(id)

	if num == 0 {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func getUserFromAuth(username string) (*model.User, error) {
	user, count := repository.GetUserByUsername(username)

	if count == 0 {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

func isUserActive(c *gin.Context, user *model.User) {
		path := c.Request.URL.Path
		
		if !slices.Contains(urlToExclude, path) {

			if user.Active == 0 {
				c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/u/locked?id=%v", user.ID.String()))
				c.Abort()
			}
		}
}

func setUserGroups(c *gin.Context) {
	casbinEnforcer, hasCE := c.Get("casbinEnforcer")
	usr := c.MustGet("user")
	var user *model.User

	if usr != nil {
		user = usr.(*model.User)
	}

	if hasCE && user != nil {
		ce := casbinEnforcer.(*casbin.Enforcer)
		groups, _ := ce.GetRolesForUser(user.Username, "")
		c.Set("user_groups", groups)
	} else {
		c.Set("user_groups", []string{})
	}
}