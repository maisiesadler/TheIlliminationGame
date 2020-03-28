package apigateway

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

var errNotLoggedIn = errors.New("No user logged in")

// GetOrCreateAuthenticatedUser creates a new UserView for the logged in user
func GetOrCreateAuthenticatedUser(ctx context.Context, request *events.APIGatewayProxyRequest) (*AuthenticatedUser, error) {

	username, ok := ParseUsername(request)
	if !ok {
		return nil, errNotLoggedIn
	}

	view, err := getOrCreateUserView(ctx, username)
	if err != nil {
		fmt.Printf("Error finding userview: %v\n", err)
		return nil, err
	}

	user := &AuthenticatedUser{
		Username: view.Username,
		ViewID:   view.ID,
	}
	return user, nil
}

// ParseUsername attempts to parse the cognito username from the Authorizer
func ParseUsername(request *events.APIGatewayProxyRequest) (string, bool) {
	if claims, ok := request.RequestContext.Authorizer["claims"]; ok {
		c := claims.(map[string]interface{})
		username, ok := c["cognito:username"]
		return username.(string), ok
	}

	return "", false
}

// ResponseSuccessful returns a 200 response for API Gateway that allows cors
func ResponseSuccessful(body string) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp.Headers["Access-Control-Allow-Headers"] = "Content-Type,Authorization,dfd-auth"
	resp.Body = body
	resp.StatusCode = 200
	return resp
}
