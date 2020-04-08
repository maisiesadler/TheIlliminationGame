package illiminationtesting

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestType struct {
	Username string
}

func TestCanAddFindOnePredicate(t *testing.T) {
	username := "Test"
	test := &TestType{
		Username: username,
	}

	coll := CreateTestCollection()
	coll.SetFindOneFilter(example)

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

func example(value interface{}, filter bson.M) bool {
	test := value.(*TestType)
	return filter["username"] == test.Username
}
