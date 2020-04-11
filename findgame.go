package theilliminationgame

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindActiveGameSetUp lets a user browse active games
func FindActiveGameSetUp(user *apigateway.AuthenticatedUser) ([]*GameSetUpSummary, error) {

	ok, coll := database.GameSetUp()
	if !ok {
		return []*GameSetUpSummary{}, errors.New("Not connected")
	}

	findOptions := options.Find()
	filter := bson.D{primitive.E{Key: "active", Value: true}}

	results, err := coll.Find(context.TODO(), filter, findOptions, &models.GameSetUp{})

	games := make([]*GameSetUpSummary, 0)

	for i := range results {
		setup := i.(*models.GameSetUp)
		gs := &GameSetUp{setup}
		if gs.playerCanJoinGame(user) {
			summary := gs.Summary(user)
			games = append(games, summary)
		}
	}

	return games, err
}
