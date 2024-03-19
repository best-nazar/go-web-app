package model

import "github.com/google/uuid"

const (
	AVATAR         string = "avatar"
	DEFAULT_AVATAR string = "assets/images/avatar.png"
)

type Image struct {
	Base
	Title   string    `gorm:"size:128;not null;" binding:"required,alphanum" form:"title" json:"title"`
	Context string    `gorm:"size:128;not null;index:idx_userid_context" binding:"required,alphanum" form:"context" json:"context"`
	Path    string    `gorm:"size:128;not null" binding:"required" form:"path" json:"path"`
	UserID  uuid.UUID `gorm:"index:idx_userid_context" binding:"required" form:"userID" json:"userID"`
}
