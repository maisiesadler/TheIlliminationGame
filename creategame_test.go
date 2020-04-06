package theilliminationgame

import (
	"testing"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCanJoinGame(t *testing.T) {
	maisie := &apigateway.AuthenticatedUser{
		Nickname: "Maisie",
		ViewID:   primitive.NewObjectID(),
	}

	jenny := &apigateway.AuthenticatedUser{
		Nickname: "Jenny",
		ViewID:   primitive.NewObjectID(),
	}

	game := Create(maisie)

	if joined := game.JoinGame(jenny); !joined {
		t.Error("Could not join game")
	}
}

func TestCanAddOptions(t *testing.T) {
	maisie := &apigateway.AuthenticatedUser{
		Nickname: "Maisie",
		ViewID:   primitive.NewObjectID(),
	}

	game := Create(maisie)

	if added := game.AddOption(maisie, "Miss Congeniality"); !added {
		t.Error("Could not add option")
	}
}

func TestCannotAddDuplicates(t *testing.T) {
	maisie := &apigateway.AuthenticatedUser{
		Nickname: "Maisie",
		ViewID:   primitive.NewObjectID(),
	}

	game := Create(maisie)

	if added := game.AddOption(maisie, "Miss Congeniality"); !added {
		t.Error("Could not add option")
	}

	if added := game.AddOption(maisie, "Miss Congeniality"); added {
		t.Error("Could add duplicate option")
	}
}
