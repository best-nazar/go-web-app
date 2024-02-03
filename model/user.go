// model/user.go

package model

type User struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	Name        string        `json:"fullName" form:"full_name" binding:"required"`
	Email       string        `json:"email" form:"email" gorm:"index" binding:"required,email"`
	Birthday    string 			`json:"birthday" form:"birthday" binding:"required"`
	Token       string        `gorm:"index"`
	CreatedAt   int64         `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   int64         `json:"updatedAt" gorm:"autoUpdateTime"`
	SuspendedAt int64         `json:"suspenedAt" form:"suspended_at"`
	Active      int16 			`json:"active" form:"active"`
	Username    string        `json:"username" form:"username" gorm:"index" binding:"required,alphanum"`
	Password    string        `json:"-" form:"password" binding:"required"`
}

type UpdateUser struct {
	ID          uint          `gorm:"-" json:"id"`
	Name        string        `json:"fullName" form:"full_name" binding:"required" gorm:"-"`
	Email       string        `json:"email" form:"email" gorm:"-" binding:"required,email"`
	Birthday    string			`json:"birthday" form:"birthday" gorm:"-" binding:"required"`
	Token       string        `gorm:"-"`
	UpdatedAt   int64         `json:"updatedAt" gorm:"autoUpdateTime, -"`
	SuspendedAt int64         `json:"suspenedAt" form:"suspended_at" gorm:"-"`
	Active      int16 			`json:"active" form:"active" gorm:"-"`
}

// Checks if the password is valid
func IsPasswordValid(user User, password string) bool {
	return user.Password == password
}
