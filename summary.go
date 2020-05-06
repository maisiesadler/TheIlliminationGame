package theilliminationgame

import (
	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/models"
)

// Summary returns a summary of the game
func (g *Game) Summary(user *apigateway.AuthenticatedUser) *GameSummary {
	remaining := []string{}
	illiminated := []string{}
	players := make([]string, len(g.db.Players))

	var lastIlliminated *LastIlliminated
	if illiminatedLen := len(g.db.Actions); illiminatedLen > 0 {
		lastIlliminated = &LastIlliminated{
			mainListIndex: g.db.Actions[illiminatedLen-1].OptionIdx,
		}
	}

	var illiminatedIdxInRemaining int
	for idx, v := range g.db.Options {
		if v.Illiminated {
			illiminated = append(illiminated, v.Name)

			if lastIlliminated.mainListIndex == idx {
				lastIlliminated.Option = v.Name
				lastIlliminated.OldIndex = illiminatedIdxInRemaining
			}
		} else {
			remaining = append(remaining, v.Name)
			illiminatedIdxInRemaining++
		}
	}

	userInGame := false
	for i, v := range g.db.Players {
		players[i] = v.Nickname
		if v.ID == user.ViewID {
			userInGame = true
		}
	}

	var status string
	var winner *Winner
	if g.db.State == models.StateRunning {
		if currentPlayer := g.playersTurn(); currentPlayer != nil {
			if currentPlayer.ID == user.ViewID {
				status = "It's your turn"
			} else {
				status = "It's " + currentPlayer.Nickname + "'s turn"
			}
		}
	} else {
		_, winner = g.checkForWinner()
		status = string(g.db.State)
	}

	actions := []*Action{}

	for _, action := range g.db.Actions {
		playerIdx := action.PlayerIdx
		optionIdx := action.OptionIdx

		if playerIdx < len(g.db.Players) && optionIdx < len(g.db.Options) {
			player := g.db.Players[playerIdx]
			option := g.db.Options[optionIdx]

			actions = append(actions, &Action{
				Player: player.Nickname,
				Option: option.Name,
				Action: action.Action,
			})
		}
	}

	return &GameSummary{
		ID:              g.db.ID,
		Remaining:       remaining,
		Illiminated:     illiminated,
		Players:         players,
		SetUpCode:       g.db.SetUpCode,
		Status:          status,
		UserInGame:      userInGame,
		Winner:          winner,
		Actions:         actions,
		LastIlliminated: lastIlliminated,
	}
}

// Summary returns a summary of the game setup
func (g *GameSetUp) Summary(user *apigateway.AuthenticatedUser) *GameSetUpSummary {
	options := make([]*SetUpOption, len(g.db.Options))
	players := make([]string, len(g.db.Players))

	for i, v := range g.db.Options {
		canEdit := *v.AddedByID == user.ViewID
		options[i] = &SetUpOption{
			Name:        v.Name,
			Description: v.Description,
			Link:        v.Link,
			AddedBy:     v.AddedByName,
			CanEdit:     canEdit,
		}
	}

	userInGame := false
	for i, v := range g.db.Players {
		players[i] = v.Nickname
		if v.ID == user.ViewID {
			userInGame = true
		}
	}

	canBeStarted := g.canBeStarted(user) == CanBeStarted

	var games []*GameSummary
	if !g.db.Active {
		games, _ = FindActiveGameForSetUp(user, g.db.Code)
	}

	return &GameSetUpSummary{
		ID:           g.db.ID,
		Code:         g.db.Code,
		Games:        games,
		Options:      options,
		Players:      players,
		UserInGame:   userInGame,
		CanBeStarted: canBeStarted,
	}
}
