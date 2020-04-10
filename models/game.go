package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// GameSetUp is the game while it is being created
type GameSetUp struct {
	ID      *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Options []string            `json:"options"`
	Players []*Player           `json:"players"`
}

// Game is the running game
type Game struct {
	ID                 *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	State              State               `json:"state"`
	CurrentPlayerIndex int                 `json:"playerIdx"`
	Options            []*Option           `json:"options"`
	Players            []*Player           `json:"players"`
}

// Option is an added option
type Option struct {
	Name        string
	Illiminated bool
}

// State is the current state of the game
type State string

const (
	StateRunning  State = "Running"
	StateFinished State = "Finished"
)

// Player represents a player
type Player struct {
	ID       primitive.ObjectID
	Nickname string
}
