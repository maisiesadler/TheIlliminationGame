package testing

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/maisiesadler/theilliminationgame/apigateway"
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
