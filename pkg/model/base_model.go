package baseModel

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID 		  uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsActive  bool           `gorm:"default:true" json:"isActive"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}