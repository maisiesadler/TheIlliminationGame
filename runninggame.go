package theilliminationgame

import (
	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/models"
)

// IlliminationResult is the result of the illiminate operation
type IlliminationResult string

const (
	Illiminated        IlliminationResult = "Illiminated"
	NotRunning                            = "Not running"
	AlreadyIlliminated                    = "Already Illiminated"
	OptionNotValid                        = "Option is not valid"
	NotYourTurn                           = "Not your turn"
	DidNotSave                            = "Did not save"
)

// Illiminate will illiminate one option and move the game on
func (g *Game) Illiminate(user *apigateway.AuthenticatedUser, option string) IlliminationResult {
	g.evaluate()

	result, idx := g.illiminate(user, option)
	if result != Illiminated {
		return result
	}

	g.db.Actions = append(g.db.Actions, &models.Action{
		Action:    "Illiminated",
		PlayerIdx: g.db.CurrentPlayerIndex,
		OptionIdx: *idx,
	})
	g.moveForward()

	if saved := g.save(); !saved {
		return DidNotSave
	}

	return Illiminated
}

// Cancel cancels a running game
func (g *Game) Cancel(user *apigateway.AuthenticatedUser) bool {
	if g.db.State != models.StateRunning {
		return false
	}

	if !g.userIsInGame(user) {
		return false
	}

	g.db.State = models.StateCancelled

	return g.save()
}

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

func (g *Game) illiminate(user *apigateway.AuthenticatedUser, option string) (IlliminationResult, *int) {
	if g.db.State != models.StateRunning {
		return NotRunning, nil
	}

	currentPlayer := g.db.Players[g.db.CurrentPlayerIndex]
	if currentPlayer.ID != user.ViewID {
		return NotYourTurn, nil
	}

	for idx, o := range g.db.Options {
		if o.Name == option {
			if o.Illiminated {
				return AlreadyIlliminated, nil
			}

			o.Illiminated = true
			return Illiminated, &idx
		}
	}

	return OptionNotValid, nil
}

func (g *Game) moveForward() {
	g.evaluate()

	if g.db.State != models.StateRunning {
		return
	}

	g.db.CurrentPlayerIndex++
	if g.db.CurrentPlayerIndex >= len(g.db.Players) {
		g.db.CurrentPlayerIndex = 0
	}
}

func (g *Game) evaluate() {
	if g.db.State == models.StateRunning {
		if gameHasFinished, _ := g.checkForWinner(); gameHasFinished {
			g.db.State = models.StateFinished
		}
	}
}

func (g *Game) checkForWinner() (bool, string) {
	var remaining *string
	for _, o := range g.db.Options {
		if !o.Illiminated {
			if remaining != nil {
				return false, ""
			}

			remaining = &o.Name
		}
	}

	if remaining != nil {
		return true, *remaining
	}

	// Something has gone wrong
	return false, ""
}
