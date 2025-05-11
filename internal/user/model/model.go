package userModel

import (
	"github.com/aragornz325/piloto-api/internal/profile/model" 
	"github.com/aragornz325/piloto-api/pkg/model"
)

type User struct {
	baseModel.BaseModel
	FirstName string                  `json:"first_name"`
	LastName  string                  `json:"last_name"`
	Email     string                  `gorm:"uniqueIndex" json:"email"`
	Driver    bool                    `json:"driver"`
	Profile   *profileModel.Profile  `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"profile,omitempty"`
}

type CreateUserInput struct {
	FirstName *string `json:"first_name" binding:"required"`
	LastName  *string `json:"last_name" binding:"required"`
	Email     *string `json:"email" binding:"required,email"`
	Driver    *bool   `json:"driver"`
}
