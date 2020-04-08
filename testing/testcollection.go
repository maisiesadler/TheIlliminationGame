package testing

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestCollection wraps a map
type TestCollection struct {
	coll             map[primitive.ObjectID]interface{}
	findOnePredicate func(interface{}, bson.M) bool
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

func CreateTestCollection() *TestCollection {
	coll := make(map[primitive.ObjectID]interface{})
	return &TestCollection{coll: coll}
}

func testCollectionIsAnICollection() {
	func(coll ICollection) {}(&TestCollection{})
}

func (coll *TestCollection) InsertOne(ctx context.Context, document interface{}) (*primitive.ObjectID, error) {
	id := primitive.NewObjectID()
	coll.coll[id] = document
	return &id, nil
}

func (coll *TestCollection) InsertOneAndFind(ctx context.Context, document interface{}, output interface{}) (interface{}, error) {
	id, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return coll.FindByID(ctx, id, output)
}

func (coll *TestCollection) DeleteByID(ctx context.Context, objID *primitive.ObjectID) error {
	delete(coll.coll, *objID)
	return nil
}

func (coll *TestCollection) UpdateByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) error {
	coll.coll[*objID] = obj
	return nil
}

func (coll *TestCollection) FindByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) (interface{}, error) {
	return coll.coll[*objID], nil
}

func (coll *TestCollection) SetFindOneFilter(predicate func(interface{}, bson.M) bool) {
	coll.findOnePredicate = predicate
}

func (coll *TestCollection) FindOne(ctx context.Context, filter interface{}, obj interface{}) (interface{}, error) {
	if coll.findOnePredicate == nil {
		return nil, errors.New("Call SetFindOneFilter")
	}

	elementMap := filter.(bson.D).Map()

	for _, v := range coll.coll {
		if coll.findOnePredicate(v, elementMap) {
			return v, nil
		}
	}

	return nil, errors.New("No matching element")
}
