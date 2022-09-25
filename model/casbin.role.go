package model

const (
	//Administrator
	ADMIN_ROLE string = "admin"
	// Not logged in
	GUEST_ROLE string = "guest"
	// Regular user
	USER_ROLE string = "friend"
)

// RBAC Role
type CasbinRole struct {
	ID     		uint   	`gorm:"primaryKey"`
	Title  		string 	`json:"title" gorm:"index" form:"title" binding:"alphanum,min=3"`
	CreatedAt 	int64  	`json:"createdAt" gorm:"autoCreateTime"`
	IsSystem 	*bool 	`json:"isSystem,omitempty" gorm:"default:false" form:"isSystem"`
}