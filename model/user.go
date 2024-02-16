// model/user.go

package model

import "github.com/best-nazar/web-app/security"
//gorm:"column:password->:false" disabled read from db
//json:”-” : ignore a field

type BaseModel struct {
	CreatedAt   int64         `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   int64         `json:"updatedAt" gorm:"autoUpdateTime"`
}

type User struct {
	BaseModel
	ID          uint          `gorm:"primaryKey" json:"id"`
	Name        string        `json:"fullName" form:"full_name" binding:"required"`
	Email       string        `json:"email" form:"email" gorm:"index" binding:"required,email"`
	Birthday    string 			`json:"birthday" form:"birthday" binding:"required"`
	Token       string        `gorm:"index"`
	SuspendedAt int64         `json:"suspenedAt" form:"suspended_at"`
	Active      int16 			`json:"active" form:"active"`
	Username    string        `json:"username" form:"username" gorm:"index" binding:"required,alphanum"`
	Password    string        `json:"-" form:"password" binding:"required"`
}

//DTO
type UserDTO struct {
	ID          uint          `gorm:"-" json:"id"`
	Name        string        `json:"fullName" form:"full_name" binding:"required" gorm:"-"`
	Username	string		  `gorm:"-" json:"username" form:"username"`
	Email       string        `json:"email" form:"email" gorm:"-" binding:"required,email"`
	Birthday    string		  `json:"birthday" form:"birthday" gorm:"-" binding:"required"`
	Token       string        `gorm:"-"`
	UpdatedAt   int64         `json:"updatedAt" gorm:"autoUpdateTime, -"`
	SuspendedAt int64         `json:"suspenedAt" form:"suspended_at" gorm:"-"`
	Active      int16 		  `json:"active" form:"active" gorm:"-"`
}

// Checks if the password is valid
func (user *User) IsPasswordValid(password string) bool {
	return user.Password == security.ComputeHmac256(password)
}

func (user *User) ConvertToDTO() *UserDTO {
	return &UserDTO{
		ID: user.ID,
		Name: user.Name,
		Username: user.Username,
		Email: user.Email,
		Birthday: user.Birthday,
		Active: user.Active,
	}
}