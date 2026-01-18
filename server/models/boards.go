package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LayoutType string

const (
	GridLayout LayoutType = "grid"
	ListLayout LayoutType = "list"
)

type Board struct {
	ID            primitive.ObjectID `json:"id" bson:"_id" validate:"required"`
	AuthorId      primitive.ObjectID `json:"authorId" bson:"authorId" validate:"required"`
	Name          string             `json:"name" bson:"name" validate:"required,min=3,max=100,alphanumspace,excludesall=<>{}[]()&|\\"`
	Description   string             `json:"description" bson:"description" validate:"omitempty,max=500,excludesall=<>{}[]()&|\\"`
	Layout        LayoutType         `json:"layout" bson:"layout" validate:"required,oneof=grid list"`
	ImageUrl      string             `json:"imageUrl" bson:"imageUrl" validate:"omitempty,url,excludesall=<>{}[]()&|\\"`
	Collaborators []Collaborator     `json:"collaborators" bson:"collaborators"`
	Keys          []Key              `json:"keys" bson:"keys"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}
