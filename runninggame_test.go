package theilliminationgame

import (
	"testing"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/maisiesadler/theilliminationgame/models"
)

func TestCanPlayGame(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	setup := Create(maisie)
	setup.JoinGame(jenny)

	setup.AddOption(maisie, "Miss Congeniality")
	setup.AddOption(jenny, "Little Princess")
	setup.AddOption(jenny, "Matilda")

	game, err := setup.Start()
	if err != nil {
		t.Errorf("Error starting game: '%v'", err)
		t.FailNow()
	}

	if game.db.CurrentPlayerIndex != 0 {
		t.Error("Game did not start with first player")
	}

	result := game.Illiminate(maisie, "Little Princess")

	if result != Illiminated {
		t.Errorf("Could not illiminate film, error: '%v'", result)
	}

	if game.db.CurrentPlayerIndex != 1 {
		t.Error("Game did not move forward")
	}

	result = game.Illiminate(jenny, "Miss Congeniality")

	if result != Illiminated {
		t.Errorf("Could not illiminate film, error: '%v'", result)
	}

	if game.db.State != models.StateFinished {
		t.Errorf("Game is not finished, actual: %v", game.db.State)
	}

	summary := game.Summary(maisie)

	if summary.Status != "Finished" {
		t.Errorf("Game Summary is not Finished, actual: %v", summary.Status)
	}

	if summary.Winner != "Matilda" {
		t.Errorf("Game Summary does not have expected winner, actual: %v", summary.Winner)
	}
}
