package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID		  uint				`gorm:"primarykey" json:"id"`
	CreatedAt time.Time			`json:"createdAt"`
	UpdatedAt time.Time  		`json:"updatedAt"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"deletedAt"`
	UUID      uuid.UUID         `gorm:"type:uuid;uinque json:"uuid"` 
}


// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	base.UUID = uuid
	return nil
}