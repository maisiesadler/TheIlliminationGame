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

func (gameSetup *GameSetUp) AddAllUserOptions(user *apigateway.AuthenticatedUser) error {

	// gameSetup.removeUserOptions(user)

	for _, o := range gameSetup.db.Options {
		if *o.AddedByID == user.ViewID {
			uo := &models.UserOption{
				UserID:      user.ViewID,
				Name:        o.Name,
				Description: o.Description,
				Link:        o.Link,
				Tags:        gameSetup.db.Tags,
				GameSetupID: *gameSetup.db.ID,
			}

			o := &UserOption{
				db: uo,
			}
			o.save(context.TODO())
		}
	}

	return nil
}

func (uo *UserOption) Remove(user *apigateway.AuthenticatedUser) error {

	ok, coll := database.UserOptions()
	if !ok {
		return errors.New("Not connected")
	}

	return coll.DeleteByID(context.Background(), uo.db.ID)
}

func (gameSetup *GameSetUp) removeUserOptions(user *apigateway.AuthenticatedUser) error {

	gameSetUpID := bson.M{"gameSetupId": gameSetup.db.ID}
	idMatch := bson.M{"userid": user.ViewID}

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
