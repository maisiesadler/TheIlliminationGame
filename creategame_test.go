package theilliminationgame

import (
	"testing"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCanJoinGame(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

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

	if len(game.db.Players) != 2 {
		t.Errorf("Not correct number of players. Expected=2, Actual=%v.", len(game.db.Players))
	}
}

func TestJoinGameDoesNotDuplicatePlayers(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

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

	if joined := game.JoinGame(jenny); joined {
		t.Error("Could not join game")
	}

	if len(game.db.Players) != 2 {
		t.Errorf("Not correct number of players. Expected=2, Actual=%v.", len(game.db.Players))
	}
}

func TestOwnerCanAddOptions(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := &apigateway.AuthenticatedUser{
		Nickname: "Maisie",
		ViewID:   primitive.NewObjectID(),
	}

	game := Create(maisie)

	if added := game.AddOption(maisie, "Miss Congeniality"); !added {
		t.Error("Could not add option")
	}
}

func TestNewPlayersCanAddOptions(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

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

	if added := game.AddOption(jenny, "Miss Congeniality"); !added {
		t.Error("Could not add option")
	}
}

func TestCannotAddOptionsIfNotInGame(t *testing.T) {
	maisie := &apigateway.AuthenticatedUser{
		Nickname: "Maisie",
		ViewID:   primitive.NewObjectID(),
	}

	jenny := &apigateway.AuthenticatedUser{
		Nickname: "Jenny",
		ViewID:   primitive.NewObjectID(),
	}

	game := Create(maisie)

	if added := game.AddOption(jenny, "Miss Congeniality"); added {
		t.Error("User who is not a player could add an option")
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
