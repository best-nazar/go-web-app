// middleware.auth.go

package main

import (
	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/gin-gonic/gin"
)

// This middleware sets whether the user is logged in or not
func setUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			auth, _ := c.Cookie("auth")
			c.Request.Header.Set("Authorization", auth)
			c.Set("is_logged_in", true) // Used for UI/Menu template (see render() in main.go)

			_, id, errt := helpers.RecoverSessionToken(token)

			if errt != nil {
				c.Set("is_logged_in", false)
			} else {
				if c.FullPath() != "" {
					repository.AddUserActivity(c.FullPath(), "path", id)
				}
			}
		} else {
			c.Set("is_logged_in", false)
			// Set guest user if he's not logged in
			if c.Request.Header.Get("Authorization") == "" {
				c.Request.SetBasicAuth(model.GUEST_ROLE, "")
			}
		}
	}
}
