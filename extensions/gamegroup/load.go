package gamegroup

import (
	"context"
	"errors"

	"github.com/maisiesadler/theilliminationgame"
	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LoadGroups(user *apigateway.AuthenticatedUser) ([]*Group, error) {
	dbgroups, err := findGroupsForUser(user)
	if err != nil {
		return nil, err
	}

	groups := make([]*Group, 0)
	for _, g := range dbgroups {
		group := loadGroup(g)
		inBson := bson.M{"$in": user.ViewID}
		idBson := bson.M{"_id": inBson}
		gamesetups, err := theilliminationgame.FindGameSetupMatchingFilter(user, &[]bson.M{idBson})
		if err != nil {
			group.GameSetUp = gamesetups
		}
		games, err := theilliminationgame.FindGamesMatchingFilter(user, &[]bson.M{idBson})
		if err != nil {
			group.Games = games
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func loadGroup(g *dbGroup) *Group {
	return &Group{
		ID:      g.ID,
		Members: g.Members,
		Name:    g.Name,
	}
}

func findGroupsForUser(user *apigateway.AuthenticatedUser) ([]*dbGroup, error) {

	ok, coll := database.GameSetUp()
	if !ok {
		return []*dbGroup{}, errors.New("Not connected")
	}

	findOptions := options.Find()

	idMatch := bson.M{"members": bson.M{"$elemMatch": bson.M{"id": user.ViewID}}}
	results, err := coll.Find(context.TODO(), idMatch, findOptions, &dbGroup{})
	if err != nil {
		return []*dbGroup{}, err
	}

	groups := make([]*dbGroup, 0)

	for i := range results {
		group := i.(*dbGroup)
		groups = append(groups, group)
	}

	return groups, err
}
