package theilliminationgame

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
)

func TestCanJoinGame(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	game := Create(maisie)

	if joined := game.JoinGame(jenny); !joined {
		t.Error("Could not join game")
	}

	if len(game.db.Players) != 2 {
		t.Errorf("Not correct number of players. Expected=2, Actual=%v.", len(game.db.Players))
	}
}

func TestCanViewIfInGame(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	game := Create(maisie)

	jensSummary := game.Summary(jenny)

	assert.False(t, jensSummary.UserInGame)

	if joined := game.JoinGame(jenny); !joined {
		t.Error("Could not join game")
	}

	jensSummary = game.Summary(jenny)

	assert.True(t, jensSummary.UserInGame)
}

func TestJoinGameDoesNotDuplicatePlayers(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	game := Create(maisie)

	if joined := game.JoinGame(jenny); !joined {
		t.Error("Could not join game")
	}

	if joined := game.JoinGame(jenny); joined {
		t.Error("Could not join game")
	}

	if len(game.db.Players) != 2 {
		t.Errorf("Not correct number of players. Expected=2, Actual=%v.", len(game.db.Players))
	}
}

func TestOwnerCanAddOptions(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")

	game := Create(maisie)

	if added := game.AddOption(maisie, "Miss Congeniality"); added != AORSuccess {
		t.Errorf("Could not add option: %v", added)
	}
}

func TestNewPlayersCanAddOptions(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	game := Create(maisie)
	game.JoinGame(jenny)

	if added := game.AddOption(jenny, "Miss Congeniality"); added != AORSuccess {
		t.Errorf("Could not add option: %v", added)
	}
}

func TestCannotAddOptionsIfNotInGame(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	game := Create(maisie)

	if added := game.AddOption(jenny, "Miss Congeniality"); added != AORUserNotInGame {
		t.Error("User who is not a player could add an option")
	}
}

func TestCannotAddDuplicates(t *testing.T) {
	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")

	game := Create(maisie)

	if added := game.AddOption(maisie, "Miss Congeniality"); added != AORSuccess {
		t.Error("Could not add option")
	}

	if added := game.AddOption(maisie, "Miss Congeniality"); added != AORAlreadyAdded {
		t.Error("Could add duplicate option")
	}

	if added := game.AddOption(maisie, " Miss Congeniality   "); added != AORAlreadyAdded {
		t.Error("Could add duplicate option with whitespace")
	}

	if added := game.AddOption(maisie, "Miss congeniality"); added != AORAlreadyAdded {
		t.Error("Could add duplicate option with different case")
	}
}

func TestStartedGameSetUpSummaryShowsGame(t *testing.T) {
	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetGameFindPredicate(testActiveGameForSetUpPredicate)

	maisie := illiminationtesting.TestUser(t, "Maisie")

	setup := Create(maisie)

	setup.AddOption(maisie, "One")
	setup.AddOption(maisie, "Two")
	setup.AddOption(maisie, "Three")

	game, startResult := setup.Start(maisie)
	assert.Equal(t, Success, startResult)

	setup, err := LoadGameSetUp(setup.db.ID)
	assert.Nil(t, err)

	summary := setup.Summary(maisie)

	assert.Equal(t, 1, len(summary.Games))
	assert.Equal(t, game.db.ID, summary.Games[0].ID)
}

func TestCanUpdateOwnOptions(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	game := Create(maisie)
	game.JoinGame(jenny)

	game.AddOption(maisie, "One")
	game.AddOption(jenny, "Two")

	updates := make(map[string]string)
	updates["name"] = "OneUpdated"
	if updated := game.UpdateOption(maisie, 0, updates); !updated {
		t.Error("Could add update option")
	}

	game, _ = LoadGameSetUp(game.db.ID)

	assert.Equal(t, "OneUpdated", game.db.Options[0].Name)

	if updated := game.UpdateOption(maisie, 1, updates); updated {
		t.Error("Could update someone elses options")
	}
}
