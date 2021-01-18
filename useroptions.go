package theilliminationgame

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddUserOption(user *apigateway.AuthenticatedUser, option string, description string, link string, gameSetupId primitive.ObjectID, tags []string) error {

	uo := &models.UserOption{
		UserID:      user.ViewID,
		Name:        option,
		Description: description,
		Link:        link,
		Tags:        tags,
		GameSetupID: gameSetupId,
	}

	o := &UserOption{
		db: uo,
	}
	o.save(context.TODO())

	return nil
}

func FindAllOptionsForUser(user *apigateway.AuthenticatedUser) ([]*UserOptionSummary, error) {
	options, err := findAllOptionsForUser(user)
	if err != nil {
		return nil, err
	}

	summaries := make(map[string]*UserOptionSummary)

	for _, option := range options {
		key := "name_" + option.Name
		if len(option.Link) > 0 {
			key = "link_" + option.Name
		}
		if _, ok := summaries[key]; !ok {
			summaries[key] = &UserOptionSummary{
				Name:        option.Name,
				Description: option.Description,
				Link:        option.Link,
			}
		}

		summaries[key].Tags = append(summaries[key].Tags, option.Tags...)
		summaries[key].GameSetupIDs = append(summaries[key].GameSetupIDs, option.GameSetupID)
	}

	summaryList := []*UserOptionSummary{}
	for _, v := range summaries {
		summaryList = append(summaryList, v)
	}

	return summaryList, nil
}

func findAllOptionsForUser(user *apigateway.AuthenticatedUser) ([]*models.UserOption, error) {

	idMatch := bson.M{"userid": user.ViewID}

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

func (game *Game) RebuildUserOptions(user *apigateway.AuthenticatedUser) error {

	gameSetup, err := LoadGameSetUp(&game.db.SetUpID)
	if err != nil {
		return err
	}

	gameSetup.removeUserOptions(user)

	for _, o := range gameSetup.db.Options {
		if *o.AddedByID == user.ViewID {
			AddUserOption(user, o.Name, o.Description, o.Link, *gameSetup.db.ID, game.db.Tags)
		}
	}

	return nil
}

func (gameSetup *GameSetUp) RebuildUserOptions(user *apigateway.AuthenticatedUser) error {

	gameSetup.removeUserOptions(user)

	for _, o := range gameSetup.db.Options {
		if *o.AddedByID == user.ViewID {
			AddUserOption(user, o.Name, o.Description, o.Link, *gameSetup.db.ID, gameSetup.db.Tags)
		}
	}

	return nil
}

func (gameSetup *GameSetUp) removeUserOptions(user *apigateway.AuthenticatedUser) error {

	gameSetUpID := bson.M{"gameSetupId": gameSetup.db.ID}
	idMatch := bson.M{"userId": user.ViewID}

	andBson := []bson.M{gameSetUpID, idMatch}

	ok, coll := database.UserOptions()
	if !ok {
		return errors.New("Not connected")
	}

	findOptions := options.Find()

	results, err := coll.Find(context.TODO(), andBson, findOptions, &models.UserOption{})
	if err != nil {
		return err
	}

	for i := range results {
		option := i.(*models.UserOption)
		coll.DeleteByID(context.Background(), option.ID)
	}

	return nil
}
