package illiminationtesting

import (
	"context"
	"errors"

	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestCollection wraps a map
type TestCollection struct {
	coll          map[primitive.ObjectID]interface{}
	findPredicate func(interface{}, bson.M) bool
}

func CreateTestCollection() *TestCollection {
	coll := make(map[primitive.ObjectID]interface{})
	return &TestCollection{coll: coll}
}

func testCollectionIsAnICollection() {
	func(coll database.ICollection) {}(&TestCollection{})
}

func (coll *TestCollection) InsertOne(ctx context.Context, document interface{}) (*primitive.ObjectID, error) {
	id := primitive.NewObjectID()
	coll.coll[id] = document

	// try update id on model
	if u, ok := document.(*models.UserView); ok {
		u.ID = &id
	}
	if s, ok := document.(*models.GameSetUp); ok {
		s.ID = &id
	}
	if g, ok := document.(*models.Game); ok {
		g.ID = &id
	}

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

func (coll *TestCollection) Find(ctx context.Context, filter interface{}, findOptions *options.FindOptions, obj interface{}) (<-chan interface{}, error) {
	if coll.findPredicate == nil {
		return nil, errors.New("Call SetFindFilter")
	}

	results := make(chan interface{})

	go func() {
		defer close(results)

		elementMap := filter.(bson.D).Map()

		for _, v := range coll.coll {
			if coll.findPredicate(v, elementMap) {
				results <- v
			}
		}
	}()

	return results, nil
}

func (coll *TestCollection) FindByID(ctx context.Context, objID *primitive.ObjectID, obj interface{}) (interface{}, error) {
	return coll.coll[*objID], nil
}

func (coll *TestCollection) SetFindFilter(predicate func(interface{}, bson.M) bool) {
	coll.findPredicate = predicate
}

func (coll *TestCollection) FindOne(ctx context.Context, filter interface{}, obj interface{}) (interface{}, error) {
	if coll.findPredicate == nil {
		return nil, errors.New("Call SetFindFilter")
	}

	var elementMap bson.M
	if d, ok := filter.(bson.D); ok {
		elementMap = d.Map()
	} else {
		elementMap = filter.(bson.M)
	}

	for _, v := range coll.coll {
		if coll.findPredicate(v, elementMap) {
			return v, nil
		}
	}

	return nil, errors.New("No matching element")
}
