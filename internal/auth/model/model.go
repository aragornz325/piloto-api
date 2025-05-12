package authModel

import "github.com/google/uuid"

type LoginDTO struct {
	Email          *string `json:"email" binding:"required"`
	Password       *string `json:"password" binding:"omitempty"`
}

type RegisterDTO struct {
	FirstName *string `json:"first_name" binding:"required"`
	LastName  *string `json:"last_name" binding:"required"`
	Email     *string `json:"email" binding:"required,email"`
	Driver    *bool   `json:"driver"`
}

type TokenPayload struct {
	UserId   uuid.UUID 	`json:"user_id"`
	Email    string 	`json:"email"`
	Role     string 	`json:"role"`
	Exp      int64  	`json:"exp"`
}