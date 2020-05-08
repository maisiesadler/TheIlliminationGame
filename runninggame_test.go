package theilliminationgame

import (
	"testing"

	"github.com/stretchr/testify/assert"

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

	assert.NotNil(t, summary.Winner)
	assert.Equal(t, "Matilda", summary.Winner.Name, "Game Summary does not have expected winner")
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
	assert.Equal(t, 0, len(maisiesSummary.Actions))
	assert.Nil(t, maisiesSummary.LastIlliminated)

	result := game.Illiminate(maisie, "Little Princess")
	assert.Equal(t, Illiminated, result)

	maisiesSummary = game.Summary(maisie)
	assert.Equal(t, 1, len(maisiesSummary.Actions))
	assert.Equal(t, "Maisie", maisiesSummary.Actions[0].Player)
	assert.Equal(t, "Little Princess", maisiesSummary.Actions[0].Option)
	assert.Equal(t, "Illiminated", maisiesSummary.Actions[0].Action)
	assert.NotNil(t, maisiesSummary.LastIlliminated)
	assert.Equal(t, 1, maisiesSummary.LastIlliminated.OldIndex)

	result = game.Illiminate(jenny, "Miss Congeniality")
	assert.Equal(t, Illiminated, result)

	maisiesSummary = game.Summary(maisie)
	assert.Equal(t, 2, len(maisiesSummary.Actions))
	assert.Equal(t, "Jenny", maisiesSummary.Actions[1].Player)
	assert.Equal(t, "Miss Congeniality", maisiesSummary.Actions[1].Option)
	assert.Equal(t, "Illiminated", maisiesSummary.Actions[1].Action)
	assert.NotNil(t, maisiesSummary.LastIlliminated)
	assert.Equal(t, 0, maisiesSummary.LastIlliminated.OldIndex)
}

func TestFinishedGameCanBeArchived(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	games, err := FindFinishedGame(maisie)
	assert.Nil(t, err)
	beforelen := len(games)

	setup := Create(maisie)

	setup.AddOption(maisie, "Miss Congeniality")
	setup.AddOption(maisie, "Little Princess")

	game, startResult := setup.Start(maisie)
	if startResult != Success {
		t.Errorf("Error starting game: '%v'", startResult)
		t.FailNow()
	}

	result := game.Illiminate(maisie, "Little Princess")

	if result != Illiminated {
		t.Errorf("Could not illiminate film, error: '%v'", result)
	}

	if game.db.State != models.StateFinished {
		t.Errorf("Game is not finished, actual: %v", game.db.State)
	}

	games, err = FindFinishedGame(maisie)
	assert.Nil(t, err)
	assert.Len(t, games, beforelen+1)

	ok := game.Archive(maisie)
	assert.True(t, ok, "Game could not be archived")

	games, err = FindFinishedGame(maisie)
	assert.Nil(t, err)
	assert.Len(t, games, beforelen)
}

func TestDescriptionAndLinkVisibleForWinner(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetGameFindWithStatePredicate()

	maisie := illiminationtesting.TestUser(t, "Maisie")

	setup := Create(maisie)

	setup.AddOption(maisie, "Miss Congeniality")
	setup.AddOption(maisie, "Little Princess")

	updates := make(map[string]string)
	updates["description"] = "description"
	updates["link"] = "link"
	ok := setup.UpdateOption(maisie, 0, updates)
	assert.True(t, ok)

	game, startResult := setup.Start(maisie)
	if startResult != Success {
		t.Errorf("Error starting game: '%v'", startResult)
		t.FailNow()
	}

	result := game.Illiminate(maisie, "Little Princess")

	if result != Illiminated {
		t.Errorf("Could not illiminate film, error: '%v'", result)
	}

	if game.db.State != models.StateFinished {
		t.Errorf("Game is not finished, actual: %v", game.db.State)
	}

	summary := game.Summary(maisie)
	assert.NotNil(t, summary.Winner)
	assert.Equal(t, "Miss Congeniality", summary.Winner.Name)
	assert.Equal(t, "description", summary.Winner.Description)
	assert.Equal(t, "link", summary.Winner.Link)
}
