package illiminationtesting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-lambda-go/events"
	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson"
)

// TestUser returns a user that exists in the database to create games with
func TestUser(t *testing.T, name string) *apigateway.AuthenticatedUser {

	request := CreateTestAuthorizedRequest("Test_" + name)
	user, err := apigateway.GetOrCreateAuthenticatedUser(context.TODO(), request)
	if err != nil {
		t.Errorf("Error creating user: '%v'", err)
		t.FailNow()
	}

	err = user.SetNickname(context.TODO(), name)
	if err != nil {
		t.Errorf("Error setting nickname: '%v'", err)
	}

	// Reload user
	user, err = apigateway.GetOrCreateAuthenticatedUser(context.TODO(), request)
	assert.Nil(t, err)

	return user
}

// CreateTestAuthorizedRequest creates an authenticated api gateway request for the given user
func CreateTestAuthorizedRequest(username string) *events.APIGatewayProxyRequest {
	claims := make(map[string]interface{})
	claims["cognito:username"] = username
	authorizer := make(map[string]interface{})
	authorizer["claims"] = claims
	context := events.APIGatewayProxyRequestContext{
		Authorizer: authorizer,
	}
	request := &events.APIGatewayProxyRequest{
		RequestContext: context,
	}

	return request
}

var overrides = make(map[string]*TestCollection)

// SetTestCollectionOverride sets a the database package to use a TestCollection
func SetTestCollectionOverride() {
	database.SetOverride(overrideDb)

	// set defaults for FindFilter
	SetUserViewFindPredicate(func(uv *models.UserView, m bson.M) bool {
		return m["username"] == uv.Username
	})

	SetGameFindWithForSetUpPredicate()
}

// SetGameFindPredicate overrides the logic to get the result for Find
func SetGameFindPredicate(predicate func(*models.Game, bson.M) bool) bool {
	fn := func(value interface{}, filter bson.M) bool {
		uv := value.(*models.Game)
		return predicate(uv, filter)
	}

	coll := getOrAddTestCollection("theilliminationgame", "games")
	coll.findPredicate = fn
	return true
}

// SetGameSetUpFindPredicate overrides the logic to get the result for Find
func SetGameSetUpFindPredicate(predicate func(*models.GameSetUp, bson.M) bool) bool {
	fn := func(value interface{}, filter bson.M) bool {
		uv := value.(*models.GameSetUp)
		return predicate(uv, filter)
	}

	coll := getOrAddTestCollection("theilliminationgame", "gamesetup")
	coll.findPredicate = fn
	return true
}

// SetUserViewFindPredicate overrides the logic to get the result for Find
func SetUserViewFindPredicate(predicate func(*models.UserView, bson.M) bool) bool {
	fn := func(value interface{}, filter bson.M) bool {
		uv := value.(*models.UserView)
		return predicate(uv, filter)
	}

	coll := getOrAddTestCollection("theilliminationgame", "users")
	coll.findPredicate = fn
	return true
}

// SetUserOptionsFindPredicate overrides the logic to get the result for Find
func SetUserOptionsFindPredicate(predicate func(*models.UserOption, bson.M) bool) bool {
	fn := func(value interface{}, filter bson.M) bool {
		uv := value.(*models.UserOption)
		return predicate(uv, filter)
	}

	coll := getOrAddTestCollection("theilliminationgame", "useroptions")
	coll.findPredicate = fn
	return true
}

func overrideDb(database string, collection string) database.ICollection {
	return getOrAddTestCollection(database, collection)
}

func getOrAddTestCollection(database string, collection string) *TestCollection {
	key := database + "_" + collection
	if val, ok := overrides[key]; ok {
		return val
	}
	overrides[key] = CreateTestCollection()
	return overrides[key]
}
