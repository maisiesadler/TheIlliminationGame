package database

import "go.mongodb.org/mongo-driver/mongo"

// MongoCollection wraps a connected mongo collection
type MongoCollection struct {
	MongoCollection *mongo.Collection
}

// CreateCollection gets a wrapped reference to a mongo collection
func CreateCollection(database string, collection string) (bool, *MongoCollection) {
	if mongoClient == nil {
		return false, nil
	}
	return true, &MongoCollection{mongoClient.Database(database).Collection(collection)}
}

// UserView returns a MongoCollection for the mongodb collection users
func UserView() (bool, *MongoCollection) {
	return CreateCollection("theilliminationgame", "users")
}
