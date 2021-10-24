package gamegroup

import (
	"github.com/maisiesadler/theilliminationgame"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type dbGroup struct {
	ID      *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string              `json:"name"`
	Members []GroupMember       `json:"members"`
}

type Group struct {
	ID        *primitive.ObjectID                     `json:"id" bson:"_id,omitempty"`
	Name      string                                  `json:"name"`
	Members   []GroupMember                           `json:"members"`
	GameSetUp []*theilliminationgame.GameSetUpSummary `json:"gameSetup"`
	Games     []*theilliminationgame.GameSummary      `json:"games"`
}

type GroupMember struct {
	ID       primitive.ObjectID `json:"id"`
	Nickname string             `json:"nickname"`
}
