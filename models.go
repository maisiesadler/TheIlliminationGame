package theilliminationgame

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Game is the
type Game struct {
	State              State
	CurrentPlayerIndex int
	Options            []*Option
	Players            []*Player
}

// Option is an added option
type Option struct {
	Name        string
	Illiminated bool
}

// State is the current state of the game
type State string

const (
	StateNotStarted State = "Not yet started"
	StateRunning    State = "Running"
	StateFinished   State = "Finished"
)

// Player represents a player
type Player struct {
	ID       primitive.ObjectID
	Nickname string
}
