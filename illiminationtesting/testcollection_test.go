package illiminationtesting

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestType struct {
	Username string
}

func TestCanAddFindPredicate(t *testing.T) {
	username := "Test"
	test := &TestType{
		Username: username,
	}

	coll := CreateTestCollection()
	coll.SetFindFilter(example)

	_, err := coll.InsertOne(context.TODO(), test)
	if err != nil {
		t.Errorf("Error inserting: %v", err)
		t.FailNow()
	}

	found, err := coll.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: username}}, &TestType{})
	if err != nil {
		t.Errorf("Error finding: %v", err)
	}

	if found == nil {
		t.Errorf("Did not find user: %v", err)
		t.FailNow()
	}

	if found.(*TestType).Username != username {
		t.Errorf("Did not have expected username: %v", found.(*TestType).Username)
	}
}

func TestCanFind(t *testing.T) {
	username := "Test"

	coll := CreateTestCollection()
	coll.SetFindFilter(example)

	_, err := coll.InsertOne(context.TODO(), &TestType{
		Username: username,
	})
	if err != nil {
		t.Errorf("Error inserting: %v", err)
		t.FailNow()
	}

	coll.InsertOne(context.TODO(), &TestType{
		Username: username,
	})
	coll.InsertOne(context.TODO(), &TestType{
		Username: "another name",
	})

	findOptions := options.Find()
	filter := bson.D{primitive.E{Key: "username", Value: username}}

	results, err := coll.Find(context.TODO(), filter, findOptions, &TestType{})
	if err != nil {
		t.Errorf("Error finding: %v", err)
	}

	parsed := make([]*TestType, 0)

	for i := range results {
		parsed = append(parsed, i.(*TestType))
	}

	if len(parsed) != 2 {
		t.Errorf("Did not get expected number of results. Expected 2, Actual %v.", len(parsed))
	}
}

func example(value interface{}, filter bson.M) bool {
	test := value.(*TestType)
	return filter["username"] == test.Username
}
