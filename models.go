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
	ID           *primitive.ObjectID `json:"id"`
	Code         string              `json:"code"`
	Options      []string            `json:"options"`
	Players      []string            `json:"players"`
	UserInGame   bool                `json:"userInGame"`
	CanBeStarted bool                `json:"canBeStarted"`
	Games        []*GameSummary      `json:"games"`
}

// GameSummary is a view of the game
type GameSummary struct {
	ID          *primitive.ObjectID `json:"id"`
	Illiminated []string            `json:"illiminated"`
	Remaining   []string            `json:"remaining"`
	Players     []string            `json:"players"`
	Status      string              `json:"status"`
	SetUpCode   string              `json:"setUpCode"`
	UserInGame  bool                `json:"userInGame"`
	Winner      string              `json:"winner"`
}
