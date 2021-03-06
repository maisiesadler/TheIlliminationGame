package theilliminationgame

import (
	"context"
	"errors"

	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoadGame creates a playable game by ID
func LoadGame(id *primitive.ObjectID) (*Game, error) {
	ok, coll := database.Game()
	if !ok {
		return nil, errors.New("Not connected")
	}
	obj, err := coll.FindByID(context.TODO(), id, &models.Game{})
	if err != nil {
		return nil, err
	}

	game := &Game{
		db: obj.(*models.Game),
	}

	if game.isFinished() && game.db.CompletedGame == nil {
		game.db.CompletedGame = createCompletedGame()
	}

	return game, nil
}

// LoadGameSetUp creates a playable game set up by ID
func LoadGameSetUp(id *primitive.ObjectID) (*GameSetUp, error) {
	ok, coll := database.GameSetUp()
	if !ok {
		return nil, errors.New("Not connected")
	}
	gameSetUp, err := coll.FindByID(context.TODO(), id, &models.GameSetUp{})
	if err != nil {
		return nil, err
	}

	return &GameSetUp{
		db: gameSetUp.(*models.GameSetUp),
	}, nil
}

// LoadUserOption creates an interactable user option
func LoadUserOption(id *primitive.ObjectID) (*UserOption, error) {
	ok, coll := database.UserOptions()
	if !ok {
		return nil, errors.New("Not connected")
	}
	uo, err := coll.FindByID(context.TODO(), id, &models.UserOption{})
	if err != nil {
		return nil, err
	}

	return &UserOption{
		db: uo.(*models.UserOption),
	}, nil
}
