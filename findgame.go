package theilliminationgame

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"

	"go.mongodb.org/mongo-driver/bson"
)

// FindActiveGameSetUp lets a user browse active games they are in
func FindActiveGameSetUp(user *apigateway.AuthenticatedUser) ([]*GameSetUpSummary, error) {
	// { players: { $elemMatch: { "nickname": "Jenny"} } }

	filter := bson.M{"active": true}
	idMatch := bson.M{"players": bson.M{"$elemMatch": bson.M{"id": user.ViewID}}}

	andBson := []bson.M{filter, idMatch}

	return findGamesMatchingFilter(user, &andBson)
}

// FindAvailableGameSetUp lets a user browse active games they are not in
func FindAvailableGameSetUp(user *apigateway.AuthenticatedUser) ([]*GameSetUpSummary, error) {
	// { players: {"$not": { $elemMatch: { "nickname": "Jenny"} } } }

	filter := bson.M{"active": true}
	idMatch := bson.M{"players": bson.M{"$not": bson.M{"$elemMatch": bson.M{"id": user.ViewID}}}}

	andBson := []bson.M{filter, idMatch}

	return findGamesMatchingFilter(user, &andBson)
}

func findGamesMatchingFilter(user *apigateway.AuthenticatedUser, filter *[]bson.M) ([]*GameSetUpSummary, error) {

	ok, coll := database.GameSetUp()
	if !ok {
		return []*GameSetUpSummary{}, errors.New("Not connected")
	}

	findOptions := options.Find()

	results, err := coll.Find(context.TODO(), bson.M{"$and": filter}, findOptions, &models.GameSetUp{})
	if err != nil {
		return []*GameSetUpSummary{}, err
	}

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
