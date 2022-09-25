// User Casbin (RBAC)
package model

type CasbinRule struct {
	ID      uint    `gorm:"primaryKey"`
	P_type	string	`json:"p-type" gorm:"index"`
	V0		string	`json:"v0" gorm:"index"`
	V1		string	`json:"v1" gorm:"index"`
	V2		string  `json:"v2"`
	V3		string	`json:"v3"`
	V4		string	`json:"v4"`
	V5		string	`json:"v5"`
}