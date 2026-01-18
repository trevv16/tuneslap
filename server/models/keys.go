package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Key struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	BoardId      primitive.ObjectID `json:"boardId" bson:"boardId"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	AudioMediaId primitive.ObjectID `json:"audioMediaId" bson:"audioMediaId"`
	AudioUrl     string             `json:"audioUrl" bson:"audioUrl" validate:"omitempty,url,excludesall=<>{}[]()&|\\"`
	ImageMediaId primitive.ObjectID `json:"imageMediaId" bson:"imageMediaId"`
	ImageUrl     string             `json:"imageUrl" bson:"imageUrl" validate:"omitempty,url,excludesall=<>{}[]()&|\\"`
	HotKey       string             `json:"hotKey" bson:"hotKey" validate:"required,min=1,max=1"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}
