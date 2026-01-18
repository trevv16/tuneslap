package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collaborator struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
