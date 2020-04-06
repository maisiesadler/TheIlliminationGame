package theilliminationgame

import "github.com/maisiesadler/theilliminationgame/apigateway"

// Start validates the inputs and sets the status to Running
func (g *Game) Start() bool {
	if g.State != StateNotStarted {
		return false
	}

	if len(g.Players) == 0 {
		return false
	}

	if len(g.Options) == 0 {
		return false
	}

	g.State = StateRunning
	return g.save()
}

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

func (g *Game) save() bool {
	return true
}

func (g *Game) illiminate(user *apigateway.AuthenticatedUser, option string) IlliminationResult {
	if g.State != StateRunning {
		return NotRunning
	}

	currentPlayer := g.Players[g.CurrentPlayerIndex]
	if currentPlayer.ID != user.ViewID {
		return NotYourTurn
	}

	for _, o := range g.Options {
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

	if g.State != StateRunning {
		return
	}

	g.CurrentPlayerIndex++
	if g.CurrentPlayerIndex >= len(g.Players) {
		g.CurrentPlayerIndex = 0
	}
}

func (g *Game) evaluate() {
	if g.State == StateRunning {
		if gameHasFinished := g.checkForWinner(); gameHasFinished {
			g.State = StateFinished
		}
	}
}

func (g *Game) checkForWinner() bool {
	remaining := 0
	for _, o := range g.Options {
		if !o.Illiminated {
			remaining++
		}
	}

	if remaining == 1 {
		return true
	}

	return false
}
