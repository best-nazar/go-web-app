//models.user.activity.go
//
// User acvitity journal
package models

type UserActivity struct {
	ID        uint   `gorm:"primaryKey"`
	Activity  string `json:"activity"`
	Data      string `json:"data"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime"`
	UserID    uint
}
