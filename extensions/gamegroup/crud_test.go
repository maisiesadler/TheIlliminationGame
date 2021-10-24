package gamegroup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maisiesadler/theilliminationgame/illiminationtesting"
)

func TestCreateGroup(t *testing.T) {

	illiminationtesting.SetTestCollectionOverride()

	ctx := context.TODO()
	maisie := illiminationtesting.TestUser(t, "Maisie")
	group, err := Create(ctx, maisie, "weds")

	assert.Nil(t, err)

	assert.Equal(t, "weds", group.Name)
}
