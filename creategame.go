package theilliminationgame

import (
	"errors"
	"strings"
	"time"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/models"
	"github.com/maisiesadler/theilliminationgame/randomcode"
)

// Create creates a new game
func Create(user *apigateway.AuthenticatedUser) *GameSetUp {

	players := []*models.Player{&models.Player{Nickname: user.Nickname, ID: user.ViewID}}
	code := randomcode.Generate()

	gameSetUp := &models.GameSetUp{
		Active:      true,
		Code:        code,
		Players:     players,
		CreatedDate: time.Now(),
	}

	gs := asGameSetUp(gameSetUp)
	gs.save()

	return gs
}

// Start validates the inputs and sets the status to Running
func (g *GameSetUp) Start() (*Game, error) {

	if len(g.db.Players) == 0 {
		return nil, errors.New("Not enough players")
	}

	optionLen := len(g.db.Options)
	if optionLen == 0 {
		return nil, errors.New("Not enough options")
	}
	options := make([]*models.Option, optionLen)
	for i, v := range g.db.Options {
		options[i] = &models.Option{
			Name: v,
		}
	}

	g.db.Active = false
	g.save()

	game := &models.Game{
		Options:     options,
		Players:     g.db.Players,
		SetUpCode:   g.db.Code,
		State:       models.StateRunning,
		CreatedDate: time.Now(),
	}

	gm := &Game{
		db: game,
	}
	gm.save()

	return gm, nil
}

// AddOption lets a player add an option if the game has not started
func (g *GameSetUp) AddOption(user *apigateway.AuthenticatedUser, option string) bool {

	if !g.userIsInGame(user) {
		return false
	}

	lowerOption := strings.ToLower(option)
	for _, o := range g.db.Options {
		if strings.ToLower(o) == lowerOption {
			return false
		}
	}

	g.db.Options = append(g.db.Options, option)

	return g.save()
}

// JoinGame returns true if the user has joined the game
func (g *GameSetUp) JoinGame(user *apigateway.AuthenticatedUser) bool {

	if g.userIsInGame(user) {
		return false
	}

	g.db.Players = append(g.db.Players, &models.Player{
		ID:       user.ViewID,
		Nickname: user.Nickname,
	})

	return g.save()
}
