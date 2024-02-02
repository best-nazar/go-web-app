// User Casbin (RBAC)
package model

const (
	GROUP_TYPE_G = "g"
	GROUP_TYPE_P = "p"
)

var ACTIONS = []string{"*", "GET", "POST"}

type CasbinRule struct {
	ID     uint   `gorm:"primaryKey"`
	P_type string `json:"p-type" gorm:"index"`
	V0     string `json:"v0" gorm:"index"`
	V1     string `json:"v1" gorm:"index"`
	V2     string `json:"v2"`
	V3     string `json:"v3"`
	V4     string `json:"v4"`
	V5     string `json:"v5"`
}

type CasbinRuleG struct {
	ID     uint   `gorm:"primaryKey"`
	P_type string `default:"g" json:"p-type" gorm:"index"`
	V0     string `json:"v0" gorm:"index"`
	V1     string `json:"v1" gorm:"index"`
	V2     string `json:"v2"`
	V3     string `json:"v3"`
	V4     string `json:"v4"`
	V5     string `json:"v5"`
}

type CasbinRuleP struct {
	ID     uint   `gorm:"primaryKey"`
	P_type string `default:"p" json:"p-type" gorm:"index"`
	V0     string `json:"group" form:"group" gorm:"index" binding:"required"`
	V1     string `json:"route" form:"route" gorm:"index" binding:"required"`
	V2     string `json:"action" form:"action" binding:"required"`
	V3     string `json:"v3"`
	V4     string `json:"v4"`
	V5     string `json:"v5"`
}


