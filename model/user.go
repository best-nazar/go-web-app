// model/user.go

package model

import (
	"database/sql"
)

type User struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	Name        string        `json:"fullName" form:"full_name" binding:"required"`
	Email       string        `json:"email" form:"email" gorm:"index" binding:"required,email"`
	Birthday    sql.NullInt64 `json:"birthday"`
	Token       string        `json:"activity" gorm:"index"`
	CreatedAt   int64         `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   int64         `json:"updatedAt" gorm:"autoUpdateTime"`
	SuspendedAt int64         `json:"suspenedAt" form:"suspended_at"`
	Active      sql.NullInt16 `json:"active"`
	Username    string        `json:"username" form:"username" gorm:"index" binding:"required,alphanum"`
	Password    string        `json:"-" form:"password" binding:"required"`
}

// Checks if the password is valid
func IsPasswordValid(user User, password string) bool {
	return user.Password == password
}
