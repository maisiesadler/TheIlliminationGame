package theilliminationgame

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

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

	game, startResult := setup.Start(maisie)
	if startResult != Success {
		t.Errorf("Error starting game: '%v'", startResult)
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

func TestIlliminatedGamesAreMovedToCorrectArray(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	setup := Create(maisie)
	setup.JoinGame(jenny)

	setup.AddOption(maisie, "Miss Congeniality")
	setup.AddOption(jenny, "Little Princess")
	setup.AddOption(jenny, "Matilda")

	game, startResult := setup.Start(maisie)
	assert.Equal(t, Success, startResult)

	assert.Equal(t, 0, game.db.CurrentPlayerIndex)

	result := game.Illiminate(maisie, "Little Princess")
	assert.Equal(t, Illiminated, result)

	summary := game.Summary(maisie)
	assert.Equal(t, 1, len(summary.Illiminated))
	assert.Equal(t, 2, len(summary.Remaining))

	assert.Equal(t, "Little Princess", summary.Illiminated[0])
}

func TestLastActionIsUpdated(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetUserViewFindPredicate(func(uv *models.UserView, m bson.M) bool {
		return m["username"] == uv.Username
	})

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	setup := Create(maisie)

	setup.JoinGame(jenny)

	setup.AddOption(maisie, "Miss Congeniality")
	setup.AddOption(jenny, "Little Princess")
	setup.AddOption(jenny, "Matilda")

	game, startResult := setup.Start(maisie)
	assert.Equal(t, Success, startResult)

	maisiesSummary := game.Summary(maisie)
	assert.Nil(t, maisiesSummary.LastAction)

	result := game.Illiminate(maisie, "Little Princess")
	assert.Equal(t, Illiminated, result)

	maisiesSummary = game.Summary(maisie)
	assert.NotNil(t, maisiesSummary.LastAction)
	assert.Equal(t, "Maisie", maisiesSummary.LastAction.Player)
	assert.Equal(t, "Little Princess", maisiesSummary.LastAction.Option)
	assert.Equal(t, "Illiminated", maisiesSummary.LastAction.Action)

	result = game.Illiminate(jenny, "Miss Congeniality")
	assert.Equal(t, Illiminated, result)

	maisiesSummary = game.Summary(maisie)
	assert.NotNil(t, maisiesSummary.LastAction)
	assert.Equal(t, "Jenny", maisiesSummary.LastAction.Player)
	assert.Equal(t, "Miss Congeniality", maisiesSummary.LastAction.Option)
	assert.Equal(t, "Illiminated", maisiesSummary.LastAction.Action)
}
