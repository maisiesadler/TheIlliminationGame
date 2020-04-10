package theilliminationgame

import (
	"github.com/maisiesadler/theilliminationgame/models"
)

// GameSetUp is the game while it is being created
type GameSetUp struct {
	db *models.GameSetUp
}

// Game is the running game
type Game struct {
	db *models.Game
}
