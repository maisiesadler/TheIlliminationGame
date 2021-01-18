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
		return m["userid"] == uo.UserID
	})

	maisie := illiminationtesting.TestUser(t, "Maisie")

	setup := Create(maisie)

	setup.AddOption(maisie, "Miss Congeniality")
	setup.AddOption(maisie, "Little Princess")

	// Act
	options, err := FindAllOptionsForUser(maisie)

	// Assert
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(options), 2)
}

// func TestGetUserOptions(t *testing.T) {

// 	// Arrange
// 	testmaisie, _ := primitive.ObjectIDFromHex("5ea6b4659c550d9390f8719c")
// 	maisie := &apigateway.AuthenticatedUser{
// 		// Nickname: view.Nickname,
// 		// Username: view.Username,
// 		ViewID: testmaisie,
// 	}

// 	uo, _ := FindAllOptionsForUser(maisie)

// 	for _, i := range uo {
// 		fmt.Printf("%v is in %v games\n", i.Name, len(i.GameSetupIDs))
// 	}
// }

// func TestBuildMyUserOptions(t *testing.T) {

// 	// Arrange
// 	testmaisie, _ := primitive.ObjectIDFromHex("5ea6b4659c550d9390f8719c")
// 	maisie := &apigateway.AuthenticatedUser{
// 		// Nickname: view.Nickname,
// 		// Username: view.Username,
// 		ViewID: testmaisie,
// 	}

// 	games, _ := FindFinishedGame(maisie)
// 	for _, i := range games {
// 		game, err := LoadGame(i.ID)
// 		assert.Nil(t, err)

// 		err = game.RebuildUserOptions(maisie)

// 		// assert.Nil(t, err)

// 		// uo, _ := FindAllOptionsForUser(maisie)

// 		// for _, i := range uo {
// 		// 	fmt.Println(i.Name)
// 		// }
// 	}
// }

// func createTestAuthorizedRequest(username string) *events.APIGatewayProxyRequest {
// 	claims := make(map[string]interface{})
// 	claims["cognito:username"] = username
// 	authorizer := make(map[string]interface{})
// 	authorizer["claims"] = claims
// 	context := events.APIGatewayProxyRequestContext{
// 		Authorizer: authorizer,
// 	}
// 	request := &events.APIGatewayProxyRequest{
// 		RequestContext: context,
// 	}

// 	return request
// }
