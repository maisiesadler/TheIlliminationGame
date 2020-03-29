package apigateway

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
		Nickname: view.Nickname,
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

// SetNickname sets the nickname of the authenticated user
func (user *AuthenticatedUser) SetNickname(ctx context.Context, nickname string) error {
	ok, collection := database.UserView()
	if !ok {
		fmt.Println("Not connected")
		return errors.New("Not connected")
	}

	ing, err := collection.FindByID(ctx, &user.ViewID, &models.UserView{})

	if err != nil {
		return err
	}

	userview := ing.(*models.UserView)
	userview.Nickname = nickname

	return collection.UpdateByID(ctx, &user.ViewID, userview)
}

func getOrCreateUserView(ctx context.Context, username string) (*models.UserView, error) {

	ok, collection := database.UserView()
	if !ok {
		fmt.Println("Not connected")
		return nil, errors.New("Not connected")
	}

	ing, err := collection.FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}}, &models.UserView{})

	if err == nil {
		return ing.(*models.UserView), nil
	}

	view := &models.UserView{
		Username: username,
	}

	ing, err = collection.InsertOneAndFind(ctx, view, &models.UserView{})

	if err == nil {
		return ing.(*models.UserView), nil
	}

	return nil, err
}
