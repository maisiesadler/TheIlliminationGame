package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoCollection wraps a connected mongo collection
type MongoCollection struct {
	MongoCollection *mongo.Collection
}

// ICollection connects to db
type ICollection interface {
	InsertOne(ctx context.Context, document interface{}) (*primitive.ObjectID, error)
	InsertOneAndFind(ctx context.Context, document interface{}, output interface{}) (interface{}, error)
	DeleteByID(ctx context.Context, objID *primitive.ObjectID) error
	UpdateByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) error
	FindByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) (interface{}, error)
	FindOne(ctx context.Context, filter interface{}, obj interface{}) (interface{}, error)
}

func mongoCollectionIsAnICollection() {
	func(coll ICollection) {}(&MongoCollection{})
}

// CreateCollection gets a wrapped reference to a mongo collection
func CreateCollection(database string, collection string) (bool, ICollection) {

	if override, ok := tryGetOverrideFor(database, collection); ok {
		return true, override
	}

	if connected := Connect(); !connected {
		return false, nil
	}

	return true, &MongoCollection{mongoClient.Database(database).Collection(collection)}
}

// UserView returns an ICollection for the mongodb collection users
func UserView() (bool, ICollection) {
	return CreateCollection("theilliminationgame", "users")
}
