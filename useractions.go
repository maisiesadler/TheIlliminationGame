package theilliminationgame

import (
	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/models"
)

func (g *Game) userIsInGame(user *apigateway.AuthenticatedUser) bool {
	return playerIsInGame(user, g.db.Players)
}

func (g *GameSetUp) userIsInGame(user *apigateway.AuthenticatedUser) bool {
	return playerIsInGame(user, g.db.Players)
}

func (g *GameSetUp) playerCanJoinGame(user *apigateway.AuthenticatedUser) bool {
	return true
}

func playerIsInGame(user *apigateway.AuthenticatedUser, players []*models.Player) bool {
	for _, player := range players {
		if player.ID == user.ViewID {
			return true
		}
	}

	return false
}

func (g *Game) playersTurn() *models.Player {
	if g.db.State != models.StateRunning {
		return nil
	}

	return g.db.Players[g.db.CurrentPlayerIndex]
}
