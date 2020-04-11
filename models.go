package theilliminationgame

import (
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GameSetUp is the game while it is being created
type GameSetUp struct {
	db *models.GameSetUp
}

// Game is the running game
type Game struct {
	db *models.Game
}

// GameSetUpSummary is a view of the game
type GameSetUpSummary struct {
	ID      *primitive.ObjectID `json:"id"`
	Code    string              `json:"code"`
	Options []string            `json:"options"`
	Players []string            `json:"players"`
}

// GameSummary is a view of the game
type GameSummary struct {
	ID      *primitive.ObjectID `json:"id"`
	Options []string            `json:"options"`
	Players []string            `json:"players"`
	Status  string              `json:"status"`
	Winner  string              `json:"winner"`
}
