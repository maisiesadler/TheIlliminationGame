package apigateway

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthenticatedUser represents the user requesting the endpoint
type AuthenticatedUser struct {
	Username string
	ViewID   primitive.ObjectID
}
