package theilliminationgame

import (
	"encoding/json"
	"testing"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/stretchr/testify/assert"
)

func TestCanMarshallGameSetUpSummary(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	setup := Create(maisie)
	setup.JoinGame(jenny)

	setup.AddOption(maisie, "Miss Congeniality")
	setup.AddOption(jenny, "Little Princess")

	summary := setup.Summary(maisie)
	assert.NotNil(t, summary)

	// Act
	b, err := json.Marshal(summary)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, b)
}

func TestCanMarshallGameSummary(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")
	id := runTestGame(t, []*apigateway.AuthenticatedUser{maisie, jenny})

	game, err := LoadGame(id)
	assert.Nil(t, err)
	assert.NotNil(t, game)
	assert.NotNil(t, game.db.CompletedGame)
	summary := game.Summary(maisie)
	assert.NotNil(t, summary)

	// Act
	b, err := json.Marshal(summary)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, b)
}
