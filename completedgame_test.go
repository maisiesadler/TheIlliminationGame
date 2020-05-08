package theilliminationgame

import (
	"testing"

	"github.com/maisiesadler/theilliminationgame/apigateway"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/maisiesadler/theilliminationgame/models"
	"github.com/stretchr/testify/assert"
)

func TestCanAddReview(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")
	id := runTestGame(t, []*apigateway.AuthenticatedUser{maisie, jenny})

	game, err := LoadGame(id)
	assert.Nil(t, err)
	assert.NotNil(t, game)
	assert.NotNil(t, game.db.CompletedGame)

	// Act
	ok := game.Review(maisie, "Really great")

	// Assert
	assert.True(t, ok)
	summary := game.Summary(maisie)
	assert.NotNil(t, summary.CompletedGame)
	assert.Len(t, summary.CompletedGame.PlayerReviews, 1)
	review, ok := summary.CompletedGame.PlayerReviews[maisie.ViewID.Hex()]
	assert.True(t, ok)
	assert.Equal(t, "Really great", review.Thoughts)
}

func TestUserCanAddOneReview(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")
	id := runTestGame(t, []*apigateway.AuthenticatedUser{maisie, jenny})

	game, err := LoadGame(id)
	assert.Nil(t, err)
	assert.NotNil(t, game)
	assert.NotNil(t, game.db.CompletedGame)

	// Act
	ok := game.Review(maisie, "Really great")
	assert.True(t, ok)
	ok = game.Review(maisie, "Actually it was rubbish")
	assert.True(t, ok)

	// Assert
	summary := game.Summary(maisie)
	assert.NotNil(t, summary.CompletedGame)
	assert.Len(t, summary.CompletedGame.PlayerReviews, 1)

	review, ok := summary.CompletedGame.PlayerReviews[maisie.ViewID.Hex()]
	assert.True(t, ok)
	assert.Equal(t, "Actually it was rubbish", review.Thoughts)
}

func TestEachUserCanAddOneReview(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")
	dad := illiminationtesting.TestUser(t, "Dad")
	id := runTestGame(t, []*apigateway.AuthenticatedUser{maisie, jenny, dad})

	game, err := LoadGame(id)
	assert.Nil(t, err)
	assert.NotNil(t, game)
	assert.NotNil(t, game.db.CompletedGame)

	// Act
	ok := game.Review(maisie, "Really great")
	assert.True(t, ok)
	ok = game.Review(maisie, "Actually it was rubbish")
	assert.True(t, ok)
	ok = game.Review(dad, "I thought it was great")
	assert.True(t, ok)
	ok = game.Review(jenny, "I fell asleep")
	assert.True(t, ok)

	// Assert
	summary := game.Summary(maisie)
	assert.NotNil(t, summary.CompletedGame)
	assert.Len(t, summary.CompletedGame.PlayerReviews, 3)

	review, ok := summary.CompletedGame.PlayerReviews[maisie.ViewID.Hex()]
	assert.True(t, ok)
	assert.Equal(t, "Actually it was rubbish", review.Thoughts)

	review, ok = summary.CompletedGame.PlayerReviews[jenny.ViewID.Hex()]
	assert.True(t, ok)
	assert.Equal(t, "I fell asleep", review.Thoughts)

	review, ok = summary.CompletedGame.PlayerReviews[dad.ViewID.Hex()]
	assert.True(t, ok)
	assert.Equal(t, "I thought it was great", review.Thoughts)
}

func TestCannotAddReviewToGameNotIn(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")
	dad := illiminationtesting.TestUser(t, "Dad")
	id := runTestGame(t, []*apigateway.AuthenticatedUser{maisie, jenny})

	game, err := LoadGame(id)
	assert.Nil(t, err)
	assert.NotNil(t, game)
	assert.NotNil(t, game.db.CompletedGame)

	// Act
	ok := game.Review(dad, "Rubbish")

	// Assert
	assert.False(t, ok)
	summary := game.Summary(maisie)
	assert.NotNil(t, summary.CompletedGame)
	assert.Len(t, summary.CompletedGame.PlayerReviews, 0)
}

func TestCannotAddReviewToGameNotFinished(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")
	jenny := illiminationtesting.TestUser(t, "Jenny")

	setup := Create(maisie)
	setup.JoinGame(jenny)

	setup.AddOption(maisie, "Miss Congeniality")
	setup.AddOption(maisie, "Little Princess")
	setup.AddOption(jenny, "Matilda")

	game, startResult := setup.Start(maisie)
	if startResult != Success {
		t.Errorf("Error starting game: '%v'", startResult)
		t.FailNow()
	}
	assert.NotNil(t, game)
	assert.Nil(t, game.db.CompletedGame)

	// Act
	ok := game.Review(maisie, "Rubbish")

	// Assert
	assert.False(t, ok)
	summary := game.Summary(maisie)
	assert.Nil(t, summary.CompletedGame)
}

func runTestGame(t *testing.T, users []*apigateway.AuthenticatedUser) *primitive.ObjectID {

	usercount := len(users)
	if usercount == 0 {
		return nil
	}

	setup := Create(users[0])

	for _, user := range users[1:] {
		setup.JoinGame(user)
	}

	i := 0
	next := func() int {
		i++
		return i % usercount
	}

	setup.AddOption(users[next()], "Miss Congeniality")
	setup.AddOption(users[next()], "Little Princess")
	setup.AddOption(users[next()], "Matilda")

	game, startResult := setup.Start(users[next()])
	if startResult != Success {
		t.Errorf("Error starting game: '%v'", startResult)
		t.FailNow()
	}

	// start with player 0
	i = 0
	assert.Equal(t, 0, game.db.CurrentPlayerIndex)

	result := game.Illiminate(users[i], "Little Princess")
	assert.Equal(t, Illiminated, result, "Could not illiminate film")
	assert.Equal(t, next(), game.db.CurrentPlayerIndex, "Game did not move forward")
	assert.Equal(t, models.StateRunning, game.db.State, "Game is not running")

	result = game.Illiminate(users[i], "Miss Congeniality")
	assert.Equal(t, Illiminated, result, "Could not illiminate film")
	assert.Equal(t, models.StateFinished, game.db.State, "Game is not finished")

	summary := game.Summary(users[next()])

	assert.Equal(t, "Finished", summary.Status, "Game Summary is not finished")
	assert.NotNil(t, summary.Winner)
	assert.Equal(t, "Matilda", summary.Winner.Name, "Game Summary does not have expected winner")

	return summary.ID
}
