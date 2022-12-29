// model/user.go

package model

import (
	"database/sql"
)

type User struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `json:"name"`
	Email       string         `json:"email" gorm:"index"`
	Birthday    sql.NullInt64  `json:"birthday"`
	Groups      string 		   `json:"groups" form:"groups"`
	Token       string         `json:"activity" gorm:"index"`
	CreatedAt   int64          `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   int64          `json:"updatedAt" gorm:"autoUpdateTime"`
	SuspendedAt int64          `json:"suspenedAt"`
	Active      sql.NullInt16  `json:"active"`
	Username    string         `json:"username" form:"username" gorm:"index"`
	Password    string         `json:"-"`
}

// Checks if the password is valid
func IsPasswordValid(user User, password string) bool {
	return user.Password == password
}
