package profileModel

import (
	"github.com/aragornz325/piloto-api/pkg/model"
	"github.com/google/uuid"
)

type Profile struct {
	baseModel.BaseModel
	UserId uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Bio          			string `json:"bio"`
	Avatar       			string `json:"avatar"`
	InstagramURL 			string `json:"instagram_url"`
	FacebookURL  			string `json:"facebook_url"`
	TwitterURL   			string `json:"twitter_url"`
	Street       			string `json:"street"`
	City         			string `json:"city"`
	State        			string `json:"state"`
	ZipCode      			string `json:"zip_code"`
	Country      			string `json:"country"`
	PhoneNumber  			string `json:"phone_number"`
	Website      			string `json:"website"`
	Whatsapp     			string `json:"whatsapp"`
}

type UserProfileDTO struct {
	UserId       *uuid.UUID `json:"user_id" binding:"required"`
	Bio          *string `json:"bio" binding:"omitempty"`
	Avatar       *string `json:"avatar" binding:"omitempty"`
	InstagramURL *string `json:"instagram_url" binding:"omitempty"`
	FacebookURL  *string `json:"facebook_url" binding:"omitempty"`
	TwitterURL   *string `json:"twitter_url" binding:"omitempty"`
	Street       *string `json:"street" binding:"omitempty"`
	City         *string `json:"city" binding:"omitempty"`
	State        *string `json:"state" binding:"omitempty"`
	ZipCode      *string `json:"zip_code" binding:"omitempty"`
	Country      *string `json:"country" binding:"omitempty"`
	PhoneNumber  *string `json:"phone_number" binding:"omitempty"`
	Website      *string `json:"website" binding:"omitempty"`
	Whatsapp     *string `json:"whatsapp" binding:"omitempty"`
}
