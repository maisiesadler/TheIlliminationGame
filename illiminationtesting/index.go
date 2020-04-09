package illiminationtesting

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson"
)

// TestUser returns a user that exists in the database to create games with
func TestUser(name string) (*apigateway.AuthenticatedUser, error) {

	request := CreateTestAuthorizedRequest("Test_" + name)
	user, err := apigateway.GetOrCreateAuthenticatedUser(context.TODO(), request)

	return user, err
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

// SetUserViewOverride sets a the database package to use a TestCollection
func SetUserViewOverride() {
	tc := CreateTestCollection()
	database.SetOverride("theilliminationgame", "users", tc)
}

// SetUserViewFindOnePredicate overrides the logic to get the result for FindOne
func SetUserViewFindOnePredicate(predicate func(*models.UserView, bson.M) bool) func(value interface{}, filter bson.M) bool {
	return func(value interface{}, filter bson.M) bool {
		uv := value.(*models.UserView)
		return predicate(uv, filter)
	}
}
