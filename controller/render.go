// render.go

package controller

import (
	"slices"

	"github.com/best-nazar/web-app/model"
	"github.com/gin-gonic/gin"
)

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func Render(c *gin.Context, data gin.H, templateName string, httpStatus int) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(httpStatus, data)
	case "application/xml":
		// Respond with XML
		c.XML(httpStatus, data)
	default:
		// Respond with HTML
		setTemplateVars(c, data)
		c.HTML(httpStatus, templateName, data)
	}
}

func setTemplateVars(c *gin.Context, data gin.H) {
	loggedInInterface, _ := c.Get("is_logged_in")
	ug, hasUG := c.Get("user_groups") 
	isAdmin := false

	if hasUG  {
		ug := ug.([]string)
		isAdmin = slices.Contains(ug, model.ADMIN_ROLE)
	}
	data["is_logged_in"] = loggedInInterface.(bool)
	data["is_admin"] = isAdmin
	data["user"] = c.MustGet("user")
}