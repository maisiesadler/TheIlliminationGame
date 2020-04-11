package theilliminationgame

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCanFindActiveGameSetUp(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetGameSetUpFindPredicate(func(gs *models.GameSetUp, m primitive.M) bool {
		return m["active"] == gs.Active
	})

	maisie := illiminationtesting.TestUser(t, "maisie")
	active, err := FindActiveGameSetUp(maisie)
	assert.Nil(t, err)
	beforeLen := len(active)

	g := Create(maisie)

	active, err = FindActiveGameSetUp(maisie)
	assert.Nil(t, err)

	assert.Equal(t, beforeLen+1, len(active))

	g.AddOption(maisie, "One")
	g.AddOption(maisie, "Two")

	g.Start()

	active, err = FindActiveGameSetUp(maisie)
	assert.Nil(t, err)

	assert.Equal(t, beforeLen, len(active))
}
