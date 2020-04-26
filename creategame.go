package theilliminationgame

import (
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

// StartResult is the result of the illiminate operation
type StartResult string

const (
	CanBeStarted     StartResult = "Can be started"
	Success          StartResult = "Success"
	NotActive        StartResult = "Not active"
	NotEnoughPlayers StartResult = "Not enough players"
	NotEnoughOptions StartResult = "Not enough options"
	UserNotInGame    StartResult = "User not in game"
)

// AddOptionResult is the result of the illiminate operation
type AddOptionResult string

const (
	AORSuccess       AddOptionResult = "Success"
	AORUserNotInGame AddOptionResult = "User not in game"
	AORAlreadyAdded  AddOptionResult = "Option already added"
	AORDidNotSave    AddOptionResult = "Did not save"
)

// Start validates the inputs and sets the status to Running
func (g *GameSetUp) Start(user *apigateway.AuthenticatedUser) (*Game, StartResult) {

	if startResult := g.canBeStarted(user); startResult != CanBeStarted {
		return nil, startResult
	}

	options := make([]*models.Option, len(g.db.Options))
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

	return gm, Success
}

// AddOption lets a player add an option if the game has not started
func (g *GameSetUp) AddOption(user *apigateway.AuthenticatedUser, option string) AddOptionResult {

	if !g.userIsInGame(user) {
		return AORUserNotInGame
	}

	option = strings.TrimSpace(option)
	lowerOption := strings.ToLower(option)
	for _, o := range g.db.Options {
		if strings.ToLower(o) == lowerOption {
			return AORAlreadyAdded
		}
	}

	g.db.Options = append(g.db.Options, option)

	if ok := g.save(); ok {
		return AORSuccess
	}

	return AORDidNotSave
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

// Deactivate will set the active flag to false
func (g *GameSetUp) Deactivate(user *apigateway.AuthenticatedUser) bool {
	if !g.userIsInGame(user) {
		return false
	}

	g.db.Active = false

	return g.save()
}

func (g *GameSetUp) canBeStarted(user *apigateway.AuthenticatedUser) StartResult {
	if !g.db.Active {
		return NotActive
	}

	if len(g.db.Players) == 0 {
		return NotEnoughPlayers
	}

	if !g.userIsInGame(user) {
		return UserNotInGame
	}

	if len(g.db.Options) < 2 {
		return NotEnoughOptions
	}

	return CanBeStarted
}
