package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        uuid.UUID `json:"ID" form:"ID" gorm:"type:char(64);primary_key"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoCreateTime"`
	// the field can have a nil value
	//used as a soft delete marker in combination with the GORM library
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := uuid.NewUUID()

	if err != nil {
		return err
	}

	base.ID = uuid
	return
}
