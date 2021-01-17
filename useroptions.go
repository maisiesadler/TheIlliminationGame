package theilliminationgame

import (
	"context"
	"errors"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddUserOption(user *apigateway.AuthenticatedUser, option string, description string, link string, tags []string) error {

	uo := &models.UserOption{
		UserID:      user.ViewID,
		Name:        option,
		Description: description,
		Link:        link,
		Tags:        tags,
	}

	o := &UserOption{
		db: uo,
	}
	o.save(context.TODO())

	return nil
}

func FindAllOptionsForUser(user *apigateway.AuthenticatedUser) ([]*models.UserOption, error) {

	idMatch := bson.M{"userId": user.ViewID}

	ok, coll := database.UserOptions()
	if !ok {
		return []*models.UserOption{}, errors.New("Not connected")
	}

	findOptions := options.Find()

	results, err := coll.Find(context.TODO(), idMatch, findOptions, &models.UserOption{})
	if err != nil {
		return []*models.UserOption{}, err
	}

	options := make([]*models.UserOption, 0)

	for i := range results {
		option := i.(*models.UserOption)
		options = append(options, option)
	}

	return options, err
}
