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
	game, err = LoadGame(id)
	assert.Nil(t, err)
	assert.NotNil(t, game)
	summary := game.Summary(maisie)
	assert.NotNil(t, summary.CompletedGame)
	assert.Len(t, summary.CompletedGame.PlayerReviews, 1)
	review := summary.CompletedGame.PlayerReviews[0]
	assert.Equal(t, "Really great", review.Thoughts)
	assert.Equal(t, "Maisie", review.PlayerNickname)
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

	review := summary.CompletedGame.PlayerReviews[0]
	assert.Equal(t, "Actually it was rubbish", review.Thoughts)
	assert.Equal(t, "Maisie", review.PlayerNickname)
	assert.True(t, summary.CompletedGame.UserHasReviewed)
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

	contains := func(reviews []PlayerReview, nickname string, thoughts string) bool {
		for _, v := range reviews {
			if v.PlayerNickname == nickname {
				assert.Equal(t, thoughts, v.Thoughts)
				return true
			}
		}

		return false
	}

	assert.True(t, contains(summary.CompletedGame.PlayerReviews, "Maisie", "Actually it was rubbish"))
	assert.True(t, contains(summary.CompletedGame.PlayerReviews, "Jenny", "I fell asleep"))
	assert.True(t, contains(summary.CompletedGame.PlayerReviews, "Dad", "I thought it was great"))
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

func TestCanAddTagFromGame(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")

	setup := Create(maisie)
	setup.AddOption(maisie, "one")
	setup.AddOption(maisie, "two")

	game, _ := setup.Start(maisie)

	if ok := game.AddTag(maisie, "pants"); !ok {
		t.Error("Could not add tag")
	}

	assert.Equal(t, 1, len(game.db.Tags))
}

func TestCanRemoveTagFromGame(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	maisie := illiminationtesting.TestUser(t, "Maisie")

	setup := Create(maisie)
	setup.AddOption(maisie, "one")
	setup.AddOption(maisie, "two")

	game, _ := setup.Start(maisie)

	if ok := game.RemoveTag(maisie, "pants"); !ok {
		t.Error("Could not remove tag")
	}

	assert.Equal(t, 0, len(game.db.Tags))
}
