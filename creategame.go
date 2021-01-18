package theilliminationgame

import (
	"context"
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
	gs.save(context.TODO())

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
	AOREmpty         AddOptionResult = "Option is empty"
)

// Start validates the inputs and sets the status to Running
func (g *GameSetUp) Start(user *apigateway.AuthenticatedUser) (*Game, StartResult) {

	if startResult := g.canBeStarted(user); startResult != CanBeStarted {
		return nil, startResult
	}

	options := make([]*models.Option, len(g.db.Options))
	for i, v := range g.db.Options {
		options[i] = &models.Option{
			Description: v.Description,
			Link:        v.Link,
			Name:        v.Name,
		}
	}

	g.db.Active = false
	g.save(context.TODO())

	game := &models.Game{
		Options:     options,
		Players:     g.db.Players,
		SetUpID:     *g.db.ID,
		SetUpCode:   g.db.Code,
		State:       models.StateRunning,
		CreatedDate: time.Now(),
		Tags:        g.db.Tags,
	}

	gm := &Game{
		db: game,
	}
	gm.save(context.TODO())

	return gm, Success
}

// AddOption lets a player add an option if the game has not started
func (g *GameSetUp) AddOption(user *apigateway.AuthenticatedUser, option string) AddOptionResult {

	return g.AddDetailedOption(user, option, "", "")
}

// AddDetailedOption lets a player add an option if the game has not started
func (g *GameSetUp) AddDetailedOption(user *apigateway.AuthenticatedUser, option string, description string, link string) AddOptionResult {

	if result := g.addDetailedOptionInner(user, option, description, link); result != AORSuccess {
		return result
	}

	AddUserOption(user, option, description, link, *g.db.ID, g.db.Tags)

	return AORSuccess
}

func (g *GameSetUp) addDetailedOptionInner(user *apigateway.AuthenticatedUser, option string, description string, link string) AddOptionResult {

	if !g.userIsInGame(user) {
		return AORUserNotInGame
	}

	if len(option) == 0 {
		return AOREmpty
	}

	option = strings.TrimSpace(option)
	lowerOption := strings.ToLower(option)
	for _, o := range g.db.Options {
		if strings.ToLower(o.Name) == lowerOption {
			return AORAlreadyAdded
		}
	}

	g.db.Options = append(g.db.Options, &models.SetUpOption{
		Name:        option,
		AddedByID:   &user.ViewID,
		AddedByName: user.Nickname,
		Description: description,
		Link:        link,
	})

	if ok := g.save(context.TODO()); ok {
		return AORSuccess
	}

	return AORDidNotSave
}

// UpdateOption allows a user to add context to their own options
func (g *GameSetUp) UpdateOption(user *apigateway.AuthenticatedUser, optionIndex int, updates map[string]string) bool {

	if !g.userIsInGame(user) {
		return false
	}

	if len(g.db.Options) <= optionIndex {
		return false
	}

	option := g.db.Options[optionIndex]

	if *option.AddedByID != user.ViewID {
		return false
	}

	if update, ok := updates["name"]; ok {
		option.Name = update
	}

	if update, ok := updates["description"]; ok {
		option.Description = update
	}

	if update, ok := updates["link"]; ok {
		option.Link = update
	}

	return g.save(context.TODO())
}

// AddTag allows a user to add a tag to a game
func (g *GameSetUp) AddTag(user *apigateway.AuthenticatedUser, tag string) bool {

	if !g.userIsInGame(user) {
		return false
	}

	if len(tag) == 0 {
		return false
	}

	tag = strings.TrimSpace(tag)
	lowerTag := strings.ToLower(tag)
	for _, t := range g.db.Tags {
		if strings.ToLower(t) == lowerTag {
			return false
		}
	}

	g.db.Tags = append(g.db.Tags, tag)

	return g.save(context.TODO())
}

// RemoveTag allows a user to remove a tag to a game
func (g *GameSetUp) RemoveTag(user *apigateway.AuthenticatedUser, tag string) bool {

	if !g.userIsInGame(user) {
		return false
	}

	tag = strings.TrimSpace(tag)
	lowerTag := strings.ToLower(tag)

	tags := []string{}
	for _, t := range g.db.Tags {
		if strings.ToLower(t) != lowerTag {
			tags = append(tags, t)
		}
	}

	g.db.Tags = tags

	return g.save(context.TODO())
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

	return g.save(context.TODO())
}

// Deactivate will set the active flag to false
func (g *GameSetUp) Deactivate(user *apigateway.AuthenticatedUser) bool {
	if !g.userIsInGame(user) {
		return false
	}

	g.db.Active = false

	return g.save(context.TODO())
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
