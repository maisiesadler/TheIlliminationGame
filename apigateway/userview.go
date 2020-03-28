package apigateway

import (
	"context"
	"errors"
	"fmt"

	"logic/database"
	"logic/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
