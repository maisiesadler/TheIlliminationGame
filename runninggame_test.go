package theilliminationgame

import (
	"testing"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCanPlayGame(t *testing.T) {
	// game := Create([]string{"Maisie", "Jenny"}, []string{"Miss Congeniality", "Little Princess", "Submarine"})

	maisie := &apigateway.AuthenticatedUser{
		Nickname: "Maisie",
		ViewID:   primitive.NewObjectID(),
	}

	jenny := &apigateway.AuthenticatedUser{
		Nickname: "Jenny",
		ViewID:   primitive.NewObjectID(),
	}

	game := Create(maisie)
	game.JoinGame(jenny)

	game.AddOption(maisie, "Miss Congeniality")
	game.AddOption(jenny, "Little Princess")
	game.AddOption(jenny, "Matilda")

	if started := game.Start(); !started {
		t.Error("Did not start")
		t.FailNow()
	}

	if game.CurrentPlayerIndex != 0 {
		t.Error("Game did not start with first player")
	}

	result := game.Illiminate(maisie, "Little Princess")

	if result != Illiminated {
		t.Errorf("Could not illiminate film, error: '%v'", result)
	}

	if game.CurrentPlayerIndex != 1 {
		t.Error("Game did not move forward")
	}

	result = game.Illiminate(jenny, "Miss Congeniality")

	if result != Illiminated {
		t.Errorf("Could not illiminate film, error: '%v'", result)
	}

	if game.State != StateFinished {
		t.Error("Game is not finished")
	}
}
