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
	Options     []*SetUpOption      `json:"options"`
	Players     []*Player           `json:"players"`
	CreatedDate time.Time           `json:"createdDate"`
}

// Game is the running game
type Game struct {
	ID                 *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SetUpID            primitive.ObjectID  `json:"setUpID"`
	SetUpCode          string              `json:"setupCode"`
	State              State               `json:"state"`
	CurrentPlayerIndex int                 `json:"playerIdx"`
	Options            []*Option           `json:"options"`
	Players            []*Player           `json:"players"`
	CreatedDate        time.Time           `json:"createdDate"`
	Actions            []*Action           `json:"actions"`
	CompletedGame      *CompletedGame      `json:"completedGame"`
}

// CompletedGame is the running game
type CompletedGame struct {
	SetUpCode     string    `json:"setupCode"`
	Players       []*Player `json:"players"`
	StartedDate   time.Time `json:"startedDate"`
	CompletedDate time.Time `json:"completedDate"`
}

// SetUpOption is an option used in the GameSetUp
type SetUpOption struct {
	Name        string
	Description string
	Link        string
	AddedByID   *primitive.ObjectID
	AddedByName string
}

// Option is an added option
type Option struct {
	Name        string
	Description string
	Illiminated bool
	Link        string
}

// Action is an action played in this game
type Action struct {
	PlayerIdx int    `json:"playerIdx"`
	OptionIdx int    `json:"optionIdx"`
	Action    string `json:"action"`
}

// State is the current state of the game
type State string

const (
	StateRunning   State = "Running"
	StateFinished  State = "Finished"
	StateArchived  State = "Archived"
	StateCancelled State = "Cancelled"
)

// Player represents a player
type Player struct {
	ID       primitive.ObjectID
	Nickname string
}
