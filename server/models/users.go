package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`
	Email          string             `json:"email" bson:"email"`
	PasswordHash   string             `json:"passwordHash" bson:"passwordHash"`
	ResetToken     *string            `json:"resetToken" bson:"resetToken"`
	ResetExpiresAt *time.Time         `json:"resetExpiresAt" bson:"resetExpiresAt"`
	ImageUrl       string             `json:"imageUrl" bson:"imageUrl"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type UpdateUserRequest struct {
	Name           string     `json:"name" bson:"name" validate:"min=3,max=100,alphanumspace,excludesall=<>{}[]()&|\\"`
	ImageUrl       string     `json:"imageUrl" bson:"imageUrl" validate:"omitempty,url,excludesall=<>{}[]()&|\\"`
	ResetToken     *string    `json:"resetToken" bson:"resetToken" validate:"omitempty,min=32,max=255"`
	ResetExpiresAt *time.Time `json:"resetExpiresAt" bson:"resetExpiresAt" validate:"omitempty"`
	PasswordHash   string     `json:"passwordHash" bson:"passwordHash" validate:"omitempty,min=6"`
}
