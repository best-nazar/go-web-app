package dto

// DTO
type UserGroup struct {
	Username string `json:"username" form:"username" binding:"required" gorm:"-"`
	Group    string `json:"group" form:"group" binding:"required" gorm:"-"`
}
