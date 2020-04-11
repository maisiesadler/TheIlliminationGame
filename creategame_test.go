package theilliminationgame

import (
	"testing"

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

	if added := game.AddOption(maisie, "Miss Congeniality"); !added {
		t.Error("Could not add option")
	}
}

func TestNewPlayersCanAddOptions(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	game := Create(maisie)
	game.JoinGame(jenny)

	if added := game.AddOption(jenny, "Miss Congeniality"); !added {
		t.Error("Could not add option")
	}
}

func TestCannotAddOptionsIfNotInGame(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	game := Create(maisie)

	if added := game.AddOption(jenny, "Miss Congeniality"); added {
		t.Error("User who is not a player could add an option")
	}
}

func TestCannotAddDuplicates(t *testing.T) {
	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")

	game := Create(maisie)

	if added := game.AddOption(maisie, "Miss Congeniality"); !added {
		t.Error("Could not add option")
	}

	if added := game.AddOption(maisie, "Miss Congeniality"); added {
		t.Error("Could add duplicate option")
	}

	if added := game.AddOption(maisie, "Miss congeniality"); added {
		t.Error("Could add duplicate option with different case")
	}
}
