package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserView represents an authenticated user
type UserView struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string              `json:"username"`
	Nickname string              `json:"nickname"`
}

// UserOption represents the users options
type UserOption struct {
	ID          *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID  `json:"userId"`
	Name        string
	Description string
	Link        string
	Tags        []string `json:"tags"`
}
