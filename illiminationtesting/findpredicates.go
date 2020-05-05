package illiminationtesting

import (
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetGameFindWithStatePredicate() {
	SetGameFindPredicate(func(g *models.Game, m primitive.M) bool {
		andval := m["$and"].(*[]bson.M)
		stateval := (*andval)[0]["state"]

		if stateval != string(g.State) {
			return false
		}

		idval := (*andval)[1]["players"].(bson.M)["$elemMatch"].(bson.M)["id"]
		id := idval.(primitive.ObjectID)

		for _, p := range g.Players {
			if p.ID == id {
				return true
			}
		}

		return false
	})

}
