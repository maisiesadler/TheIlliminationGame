package theilliminationgame

import (
	"github.com/maisiesadler/theilliminationgame/apigateway"
)

func (g *Game) userIsInGame(user *apigateway.AuthenticatedUser) bool {
	for _, player := range g.db.Players {
		if player.ID == user.ViewID {
			return true
		}
	}

	return false
}

func (g *GameSetUp) userIsInGame(user *apigateway.AuthenticatedUser) bool {
	for _, player := range g.db.Players {
		if player.ID == user.ViewID {
			return true
		}
	}

	return false
}
