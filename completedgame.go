package theilliminationgame

import (
	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/models"
)

// Archive sets a finished game state to archived
func (g *Game) Archive(user *apigateway.AuthenticatedUser) bool {
	if g.db.State != models.StateFinished {
		return false
	}

	if !g.userIsInGame(user) {
		return false
	}

	g.db.State = models.StateArchived

	return g.save()
}

// Review allows a player to review the game
func (g *Game) Review(user *apigateway.AuthenticatedUser, thoughts string) bool {
	if g.db.State != models.StateFinished {
		return false
	}

	if !g.userIsInGame(user) {
		return false
	}

	if g.db.CompletedGame == nil {
		g.db.CompletedGame = createCompletedGame(g)
	}

	if _, ok := g.db.CompletedGame.PlayerReviews[user.ViewID]; !ok {
		g.db.CompletedGame.PlayerReviews[user.ViewID] = &models.PlayerReview{}
	}
	g.db.CompletedGame.PlayerReviews[user.ViewID].Thoughts = thoughts

	return g.save()
}
