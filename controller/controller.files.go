package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/best-nazar/web-app/service"
	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	var (
		path string
		err error
	)

	id := c.Param("id")
	user, num := repository.FindUserById(id)

	if num == 0 {
		c.AbortWithError(http.StatusNotFound, errors.New("user not found"))
	} else {
		if user.Avatar() == nil {
			path = model.DEFAULT_AVATAR
		} else {
			path = user.Avatar().Path
		}
	}

	loggedUser := c.MustGet("user").(*model.User)

	if loggedUser.ID != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if c.Request.Method == http.MethodPost {
		imageLoader := new(service.ImageLoader)
		imageLoader.LoadDefaults(c)
		path, err = imageLoader.SaveFile(c)

		if err == nil {
			img := model.Image{
				Title: "User's avatar",
				Context: model.AVATAR,
				Path: path,
				UserID: user.ID,
			}

			repository.DeleteImage(user.Avatar())
			errm := repository.SaveImage(&img)

			if errm != nil {
				c.Error(errm)
			}
			
			if len(c.Errors) == 0 {
				c.Redirect(http.StatusFound, "/member/avatar/" + fmt.Sprintf("%v", user.ID))
			}
		} else {
			c.Error(err)
		}
	}

	Render(c, gin.H{
		"title": "User's Avatar",
		"description": "Image uploader",
		"user_id": id,
		"image_path": path,
		"errors": helpers.Errors(c),
	}, "user-avatar.html", http.StatusOK)
}