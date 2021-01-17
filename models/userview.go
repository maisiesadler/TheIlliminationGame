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

// UserOptions represents the users options
type UserOptions struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID   primitive.ObjectID  `json:"userId"`
	Nickname string              `json:"nickname"`
	Option   UserOption          `json:"option"`
}

// UserOption is an option used in UserOptions
type UserOption struct {
	Name        string
	Description string
	Link        string
}
