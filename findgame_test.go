package theilliminationgame

import (
	"context"
	"testing"

	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
)

func TestCanFindActiveGameSetUp(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetGameSetUpFindPredicate(func(gs *models.GameSetUp, m primitive.M) bool {
		andval := m["$and"].(*[]bson.M)
		activeval := (*andval)[0]["active"]

		if activeval != gs.Active {
			return false
		}

		idval := (*andval)[1]["players"].(bson.M)["$elemMatch"].(bson.M)["id"]
		id := idval.(primitive.ObjectID)

		for _, p := range gs.Players {
			if p.ID == id {
				return true
			}
		}

		return false
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

	game, err := g.Start()
	assert.Nil(t, err)

	active, err = FindActiveGameSetUp(maisie)
	assert.Nil(t, err)

	assert.Equal(t, beforeLen, len(active))

	// Cleanup
	ok, coll := database.GameSetUp()
	assert.True(t, ok)

	coll.DeleteByID(context.TODO(), g.Summary(maisie).ID)

	ok, coll = database.GameSetUp()
	assert.True(t, ok)

	coll.DeleteByID(context.TODO(), game.Summary(maisie).ID)
}

func TestAnotherUserCanFindAvailableGameSetUp(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()
	illiminationtesting.SetGameSetUpFindPredicate(func(gs *models.GameSetUp, m primitive.M) bool {
		andval := m["$and"].(*[]bson.M)
		activeval := (*andval)[0]["active"]

		if activeval != gs.Active {
			return false
		}

		idval := (*andval)[1]["players"].(bson.M)["$not"].(bson.M)["$elemMatch"].(bson.M)["id"]
		id := idval.(primitive.ObjectID)

		for _, p := range gs.Players {
			if p.ID == id {
				return false
			}
		}

		return true
	})

	maisie := illiminationtesting.TestUser(t, "maisie")
	jenny := illiminationtesting.TestUser(t, "jenny")

	active, err := FindAvailableGameSetUp(jenny)
	assert.Nil(t, err)
	beforeLen := len(active)

	g := Create(maisie)

	active, err = FindAvailableGameSetUp(jenny)
	assert.Nil(t, err)

	assert.Equal(t, beforeLen+1, len(active))

	ok := g.JoinGame(jenny)
	assert.True(t, ok)

	active, err = FindAvailableGameSetUp(jenny)
	assert.Nil(t, err)

	assert.Equal(t, beforeLen, len(active))

	// Cleanup
	ok, coll := database.GameSetUp()
	assert.True(t, ok)

	coll.DeleteByID(context.TODO(), g.Summary(maisie).ID)
}
