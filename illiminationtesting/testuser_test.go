package illiminationtesting

import (
	"testing"

	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/stretchr/testify/assert"
)

func TestUserHasNickname(t *testing.T) {
	SetTestCollectionOverride()
	SetUserViewFindPredicate(func(uv *models.UserView, m bson.M) bool {
		return m["username"] == uv.Username
	})

	user := TestUser(t, "User")
	assert.Equal(t, "User", user.Nickname)
}
