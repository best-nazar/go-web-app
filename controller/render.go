// render.go

package controller

import (
	"github.com/best-nazar/web-app/helpers"
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

	data["user"] = c.MustGet("user")

	if ug, hasUG := c.Get("user_groups"); hasUG  {
		data["user_groups"] = ug.([]string)
	}

	data["is_logged_in"] = loggedInInterface.(bool)
	data["errors"] = helpers.Errors(c)
}