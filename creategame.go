package theilliminationgame

import "github.com/maisiesadler/theilliminationgame/apigateway"

// Create creates a new game
func Create(user *apigateway.AuthenticatedUser) *Game {

	players := []*Player{&Player{Nickname: user.Nickname, ID: user.ViewID}}

	return &Game{
		Players: players,
		State:   StateNotStarted,
	}
}

// AddOption lets a player add an option if the game has not started
func (g *Game) AddOption(user *apigateway.AuthenticatedUser, option string) bool {
	if g.State != StateNotStarted {
		return false
	}

	if !g.userIsInGame(user) {
		return false
	}

	for _, o := range g.Options {
		if o.Name == option {
			return false
		}
	}

	g.Options = append(g.Options, &Option{
		Name: option,
	})

	return g.save()
}

// JoinGame returns true if the user has joined the game
func (g *Game) JoinGame(user *apigateway.AuthenticatedUser) bool {
	if g.State != StateNotStarted {
		return false
	}

	if g.userIsInGame(user) {
		return false
	}

	g.Players = append(g.Players, &Player{
		ID:       user.ViewID,
		Nickname: user.Nickname,
	})

	return g.save()
}
