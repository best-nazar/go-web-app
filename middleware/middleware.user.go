package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/best-nazar/web-app/repository"
	"github.com/gin-gonic/gin"
)

var urlToExclude = []string {
	"/u/locked",
	"/u/logout",
}

func IsUserActive() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt("user_id")
		path := c.Request.URL.Path
		
		if userId != 0 && !slices.Contains(urlToExclude, path) {
			user, num := repository.FindUserById(fmt.Sprintf("%v", userId))

			if num > 0 && user.Active == 0 {
				c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/u/locked?id=%v", userId))
				c.Abort()
			}
		}

		c.Next()
	}
}