package dto

type UpdateUserDto struct {
	ID       string `json:"ID" form:"ID" binding:"uuid"`
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"  binding:"email"`
	Birthday string `json:"birthday" form:"birthday"`
}

type ActivateUserDto struct {
	ID     string `json:"ID" form:"ID" binding:"uuid"`
	Active int16  `json:"active"`
}
