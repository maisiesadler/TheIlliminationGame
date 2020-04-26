package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GameSetUp is the game while it is being created
type GameSetUp struct {
	ID          *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Active      bool                `json:"active"`
	Code        string              `json:"code"`
	Options     []string            `json:"options"`
	Players     []*Player           `json:"players"`
	CreatedDate time.Time           `json:"createdDate"`
}

// Game is the running game
type Game struct {
	ID                 *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SetUpCode          string              `json:"setupCode"`
	State              State               `json:"state"`
	CurrentPlayerIndex int                 `json:"playerIdx"`
	Options            []*Option           `json:"options"`
	Players            []*Player           `json:"players"`
	CreatedDate        time.Time           `json:"createdDate"`
	LastAction         *LastAction         `json:"lastAction"`
}

// Option is an added option
type Option struct {
	Name        string
	Illiminated bool
}

// LastAction is the last played action of this game
type LastAction struct {
	PlayerIdx int    `json:"playerIdx"`
	OptionIdx int    `json:"optionIdx"`
	Action    string `json:"action"`
}

// State is the current state of the game
type State string

const (
	StateRunning   State = "Running"
	StateFinished  State = "Finished"
	StateCancelled State = "Cancelled"
)

// Player represents a player
type Player struct {
	ID       primitive.ObjectID
	Nickname string
}
