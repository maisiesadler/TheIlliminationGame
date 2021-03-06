package theilliminationgame

import (
	"time"

	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GameSetUp is the game while it is being created
type GameSetUp struct {
	db *models.GameSetUp
}

// Game is the running game
type Game struct {
	db *models.Game
}

// UserOption is a user option
type UserOption struct {
	db *models.UserOption
}

// UserOptionSummary is a view of the users options
type UserOptionSummary struct {
	ID           primitive.ObjectID   `json:"id"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Link         string               `json:"link"`
	Tags         []string             `json:"tags"`
	GameSetupIDs []primitive.ObjectID `json:"gameSetupIds"`
}

// GameSetUpSummary is a view of the game
type GameSetUpSummary struct {
	ID           *primitive.ObjectID `json:"id"`
	Code         string              `json:"code"`
	Options      []*SetUpOption      `json:"options"`
	Players      []string            `json:"players"`
	UserInGame   bool                `json:"userInGame"`
	CanBeStarted bool                `json:"canBeStarted"`
	Games        []*GameSummary      `json:"games"`
	Tags         []string            `json:"tags"`
}

// GameSummary is a view of the game
type GameSummary struct {
	ID              *primitive.ObjectID `json:"id"`
	Illiminated     []string            `json:"illiminated"`
	Remaining       []string            `json:"remaining"`
	Players         []string            `json:"players"`
	Status          string              `json:"status"`
	SetUpCode       string              `json:"setUpCode"`
	UserInGame      bool                `json:"userInGame"`
	Winner          *Winner             `json:"winner"`
	Actions         []*Action           `json:"actions"`
	LastIlliminated *LastIlliminated    `json:"lastIlliminated"`
	CompletedGame   *CompletedGame      `json:"completedGame"`
	StartedDate     time.Time           `json:"startedDate"`
	Tags            []string            `json:"tags"`
}

// CompletedGame is extra info for a completed game
type CompletedGame struct {
	CompletedDate   time.Time      `json:"completedDate"`
	PlayerReviews   []PlayerReview `json:"playerReviews"`
	UserHasReviewed bool           `json:"userHasReviewed"`
}

// PlayerReview is a users thoughts on a completed game
type PlayerReview struct {
	PlayerNickname string  `json:"playerNickname"`
	Thoughts       string  `json:"thoughts"`
	ImageURL       *string `json:"imageUrl"`
	IsUsersReview  bool    `json:"isUsersReview"`
}

// Action is an action played in the game
type Action struct {
	Player string `json:"player"`
	Option string `json:"option"`
	Action string `json:"action"`
	Link   string `json:"link"`
}

// LastIlliminated is the last illiminated option
type LastIlliminated struct {
	Option        string `json:"option"`
	OldIndex      int    `json:"oldIndex"`
	mainListIndex int
}

// Winner is info about the winning option
type Winner struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

// SetUpOption is an option used in the GameSetUp
type SetUpOption struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	AddedBy     string `json:"addedBy"`
	CanEdit     bool   `json:"canEdit"`
}
