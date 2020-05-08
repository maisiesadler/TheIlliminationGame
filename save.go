package theilliminationgame

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/maisiesadler/theilliminationgame/database"
)

func (g *Game) save(ctx context.Context) bool {

	ok, coll := database.Game()
	if !ok {
		return false
	}

	game := g.db

	id, err := addOrUpdate(ctx, game.ID, game, coll)
	if err == nil {
		game.ID = id
		return true
	}
	fmt.Printf("Error saving game: '%v'", err)
	return false
}

func (g *GameSetUp) save(ctx context.Context) bool {

	ok, coll := database.GameSetUp()
	if !ok {
		return false
	}

	setup := g.db

	id, err := addOrUpdate(ctx, setup.ID, setup, coll)
	if err == nil {
		setup.ID = id
		return true
	}
	fmt.Printf("Error saving setup: '%v'", err)
	return false
}

func addOrUpdate(ctx context.Context, id *primitive.ObjectID, doc interface{}, coll database.ICollection) (*primitive.ObjectID, error) {
	if id == nil {
		ID, err := coll.InsertOne(ctx, doc)
		return ID, err
	}

	err := coll.UpdateByID(ctx, id, doc)
	return id, err
}
