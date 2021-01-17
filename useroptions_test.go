package theilliminationgame

import (
	"testing"

	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/stretchr/testify/assert"
)

func TestUserOptionIsAddedToTable(t *testing.T) {

	// Arrange
	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetUserViewFindPredicate(func(uv *models.UserView, m primitive.M) bool {
		return m["username"] == uv.Username
	})
	illiminationtesting.SetUserOptionsFindPredicate(func(uo *models.UserOption, m primitive.M) bool {
		return m["userId"] == uo.UserID
	})

	maisie := illiminationtesting.TestUser(t, "Maisie")

	setup := Create(maisie)

	setup.AddOption(maisie, "Miss Congeniality")
	setup.AddOption(maisie, "Little Princess")

	// Act
	options, err := FindAllOptionsForUser(maisie)

	// Assert
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, 2, len(options))
}
