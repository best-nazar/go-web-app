// models/user.go

package models

import (
	"database/sql"
)

type User struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `json:"name"`
	Email       string         `json:"email" gorm:"index"`
	Birthday    sql.NullInt64  `json:"birthhday"`
	Groups      sql.NullString `json:"groups"`
	Token       string         `json:"activity" gorm:"index"`
	CreatedAt   int64          `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   int64          `json:"updatedAt" gorm:"autoUpdateTime"`
	SuspendedAt int64          `json:"suspenedAt"`
	Active      sql.NullInt16  `json:"active"`
	Username    string         `json:"username" gorm:"index"`
	Password    string         `json:"-"`
}

// Checks if the password is valid
func IsPasswordValid(user User, password string) bool {
	return user.Password == password
}
