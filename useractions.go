package theilliminationgame

import (
	"github.com/maisiesadler/theilliminationgame/apigateway"
)

func (g *Game) userIsInGame(user *apigateway.AuthenticatedUser) bool {
	for _, player := range g.Players {
		if player.ID == user.ViewID {
			return true
		}
	}

	return false
}
