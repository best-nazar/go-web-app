// model/user.go

package model

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/best-nazar/web-app/security"
)

//gorm:"column:password->:false" disabled read from db
//json:”-” : ignore a field

type User struct {
	Base
	Name        string    `json:"name" form:"name" binding:"required,min=3"`
	Email       string    `json:"email" form:"email" gorm:"index" binding:"required,email"`
	Birthday    string    `json:"birthday" form:"birthday" binding:"required"`
	Token       string    `gorm:"index"`
	SuspendedAt *time.Time `json:"suspenedAt"`
	Active      int16     `json:"active"`
	Username    string    `json:"username" form:"username" gorm:"index" binding:"required,alphanum,min=3"`
	Password    string    `json:"-" form:"password" binding:"required,min=6"`
	Images      []Image
}

func (user *User) ComparePassword(pswd string) (error, bool) {
	if user.Password == pswd {
		return nil, true
	}

	return fmt.Errorf("provided passwords does not match"), false
}

func (user *User) ValidateUsername(content string) (error, bool) {
	if slices.Contains(strings.Split(content, ","), user.Username) {
		return fmt.Errorf("username '%s' isn't available", user.Username), false
	}

	return nil, true
}

// Checks if the password is valid
func (user *User) IsPasswordValid(password string) bool {
	return user.Password == security.ComputeHmac256(password)
}

func (user *User) Avatar() *Image {
	for key, val := range user.Images {
		if val.Context == AVATAR {
			return &user.Images[key]
		}
	}

	return nil
}
