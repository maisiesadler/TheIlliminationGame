package theilliminationgame

import "github.com/maisiesadler/theilliminationgame/models"

// LoadGame creates a playable game from the db type
func LoadGame(game *models.Game) *Game {
	return &Game{
		db: game,
	}
}

// LoadGameSetUp creates a playable game set up from the db type
func LoadGameSetUp(gameSetUp *models.GameSetUp) *GameSetUp {
	return &GameSetUp{
		db: gameSetUp,
	}
}
