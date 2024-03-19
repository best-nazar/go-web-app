// model.user.activity.go
//
// User acvitity journal
package model

type UserActivity struct {
	ID        uint   `gorm:"primaryKey"`
	Activity  string `json:"activity"`
	Data      string `json:"data"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime"`
	UserID    string
}

type UsersList struct {
    Users []string `form:"users[]" json:"users[]"`
	Group string `form:"group" json:"group"`
}