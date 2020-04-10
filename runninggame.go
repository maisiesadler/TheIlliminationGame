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
	if result := g.illiminate(user, option); result != Illiminated {
		return result
	}

	g.moveForward()

	if saved := g.save(); !saved {
		return DidNotSave
	}

	return Illiminated
}

func (g *Game) illiminate(user *apigateway.AuthenticatedUser, option string) IlliminationResult {
	if g.db.State != models.StateRunning {
		return NotRunning
	}

	currentPlayer := g.db.Players[g.db.CurrentPlayerIndex]
	if currentPlayer.ID != user.ViewID {
		return NotYourTurn
	}

	for _, o := range g.db.Options {
		if o.Name == option {
			if o.Illiminated {
				return AlreadyIlliminated
			}

			o.Illiminated = true
			return Illiminated
		}
	}

	return OptionNotValid
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
		if gameHasFinished := g.checkForWinner(); gameHasFinished {
			g.db.State = models.StateFinished
		}
	}
}

func (g *Game) checkForWinner() bool {
	remaining := 0
	for _, o := range g.db.Options {
		if !o.Illiminated {
			remaining++
		}
	}

	if remaining == 1 {
		return true
	}

	return false
}
