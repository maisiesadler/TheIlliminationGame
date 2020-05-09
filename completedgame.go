package theilliminationgame

import (
	"context"

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

	return g.save(context.TODO())
}

// Review allows a player to review the game
func (g *Game) Review(user *apigateway.AuthenticatedUser, thoughts string) bool {
	if g.db.State != models.StateFinished {
		return false
	}

	if !g.userIsInGame(user) {
		return false
	}

	viewID := user.ViewID.Hex()
	if _, ok := g.db.CompletedGame.PlayerReviews[viewID]; !ok {
		g.db.CompletedGame.PlayerReviews[viewID] = &models.PlayerReview{
			PlayerNickname: user.Nickname,
		}
	}
	g.db.CompletedGame.PlayerReviews[viewID].Thoughts = thoughts

	return g.save(context.TODO())
}

// UpdateHasImage allows a player to add an image to their review
func (g *Game) UpdateHasImage(user *apigateway.AuthenticatedUser, hasImage bool) bool {
	if g.db.State != models.StateFinished {
		return false
	}

	if !g.userIsInGame(user) {
		return false
	}

	viewID := user.ViewID.Hex()
	if _, ok := g.db.CompletedGame.PlayerReviews[viewID]; !ok {
		g.db.CompletedGame.PlayerReviews[viewID] = &models.PlayerReview{
			PlayerNickname: user.Nickname,
		}
	}
	g.db.CompletedGame.PlayerReviews[viewID].Image = hasImage

	return g.save(context.TODO())
}
