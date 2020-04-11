package theilliminationgame

import (
	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/models"
)

// Summary returns a summary of the game
func (g *Game) Summary(user *apigateway.AuthenticatedUser) *GameSummary {
	options := make([]string, len(g.db.Options))
	players := make([]string, len(g.db.Players))

	for i, v := range g.db.Options {
		options[i] = v.Name
	}

	for i, v := range g.db.Players {
		players[i] = v.Nickname
	}

	var status string
	if g.db.State == models.StateRunning {
		if currentPlayer := g.playersTurn(); currentPlayer != nil {
			if currentPlayer.ID == user.ViewID {
				status = "It's your turn"
			} else {
				status = "It's " + currentPlayer.Nickname + "'s turn"
			}
		}
	} else {
		_, winner := g.checkForWinner()
		status = "Finished, winner: '" + winner + "'"
	}

	return &GameSummary{
		ID:      g.db.ID,
		Options: options,
		Players: players,
		Status:  status,
	}
}

// Summary returns a summary of the game setup
func (g *GameSetUp) Summary(user *apigateway.AuthenticatedUser) *GameSetUpSummary {
	options := make([]string, len(g.db.Options))
	players := make([]string, len(g.db.Players))

	for i, v := range g.db.Options {
		options[i] = v
	}

	for i, v := range g.db.Players {
		players[i] = v.Nickname
	}

	return &GameSetUpSummary{
		ID:      g.db.ID,
		Code:    g.db.Code,
		Options: options,
		Players: players,
	}
}
