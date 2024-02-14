// render.go

package controller

import (
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
	data["is_logged_in"] = loggedInInterface.(bool)
}