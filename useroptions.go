package theilliminationgame

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllOptionsForUser(user *apigateway.AuthenticatedUser) ([]*UserOptionSummary, error) {
	options, err := findAllOptionsForUser(user)
	if err != nil {
		return nil, err
	}

	summaries := []*UserOptionSummary{}

	for _, option := range options {
		summary := &UserOptionSummary{
			ID:           *option.ID,
			Name:         option.Name,
			Description:  option.Description,
			Link:         option.Link,
			Tags:         option.Tags,
			GameSetupIDs: option.GameSetupIDs,
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (uo *UserOption) UpdateUserOption(user *apigateway.AuthenticatedUser, updates map[string]string) error {

	if uo.db.UserID != user.ViewID {
		return errors.New("Cannot update another users options")
	}

	fmt.Println("Updating user options")

	for k, v := range updates {
		fmt.Println("Updating " + k)

		if k == "name" {
			uo.db.Name = v
		} else if k == "link" {
			uo.db.Link = v
		} else if k == "description" {
			uo.db.Description = v
		} else if k == "tag_add" {
			v = strings.ToLower(v)
			if !contains(uo.db.Tags, v) {
				uo.db.Tags = append(uo.db.Tags, v)
			}
		} else if k == "tag_remove" {
			uo.db.Tags = remove(uo.db.Tags, v)
		}
	}

	if uo.save(context.Background()) {
		return nil
	}

	return errors.New("Did not save changes")
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}

	return false
}

func remove(list []string, s string) []string {
	s = strings.ToLower(s)
	ret := []string{}
	for _, v := range list {
		if s != v {
			ret = append(ret, v)
		}
	}

	return ret
}

// func DeleteMatchingOptions(user *apigateway.AuthenticatedUser, groupKey string, gameSetupIds []primitive.ObjectID) (int, error) {

// 	ok, coll := database.UserOptions()
// 	if !ok {
// 		return 0, errors.New("Not connected")
// 	}

// 	nameOrLink, err := getGroupKeyBson(groupKey)
// 	if err != nil {
// 		return 0, err
// 	}
// 	userMatch := bson.M{"userid": user.ViewID}
// 	gameSetupID := bson.M{"gamesetupid": bson.M{"$in": gameSetupIds}}

// 	filter := []bson.M{userMatch, nameOrLink, gameSetupID}

// 	findOptions := options.Find()

// 	results, err := coll.Find(context.TODO(), bson.M{"$and": filter}, findOptions, &models.UserOption{})
// 	if err != nil {
// 		return 0, err
// 	}

// 	deleteCount := 0
// 	for i := range results {
// 		option := i.(*models.UserOption)
// 		err = coll.DeleteByID(context.Background(), option.ID)
// 		if err == nil {
// 			deleteCount++
// 		}
// 	}

// 	return deleteCount, nil
// }

// func getGroupKeyBson(groupKey string) (bson.M, error) {

// 	nameMatches := regexp.MustCompile("name_(.*)").FindStringSubmatch(groupKey)

// 	if nameMatches != nil {
// 		return bson.M{"name": nameMatches[1]}, nil
// 	}

// 	linkMatches := regexp.MustCompile("link_(.*)").FindStringSubmatch(groupKey)
// 	if linkMatches != nil {
// 		return bson.M{"link": linkMatches[1]}, nil
// 	}

// 	return nil, errors.New("Did not understand group key")
// }

func findAllOptionsForUser(user *apigateway.AuthenticatedUser) ([]*models.UserOption, error) {

	idMatch := bson.M{"userid": user.ViewID}

	ok, coll := database.UserOptions()
	if !ok {
		return []*models.UserOption{}, errors.New("Not connected")
	}

	findOptions := options.Find()

	results, err := coll.Find(context.TODO(), idMatch, findOptions, &models.UserOption{})
	if err != nil {
		return []*models.UserOption{}, err
	}

	options := make([]*models.UserOption, 0)

	for i := range results {
		option := i.(*models.UserOption)
		options = append(options, option)
	}

	return options, err
}

func (gameSetup *GameSetUp) AddAllUserOptions(user *apigateway.AuthenticatedUser) error {

	// gameSetup.removeUserOptions(user)

	for _, o := range gameSetup.db.Options {
		if *o.AddedByID == user.ViewID {

			group, _ := findGroupForUserOption(user, o.Name, o.Link)

			if group == nil {
				group = &models.UserOption{
					UserID:      user.ViewID,
					Name:        o.Name,
					Description: o.Description,
					Link:        o.Link,
				}
			}

			group.Tags = append(group.Tags, gameSetup.db.Tags...)
			group.GameSetupIDs = append(group.GameSetupIDs, *gameSetup.db.ID)

			o := &UserOption{
				db: group,
			}
			o.save(context.TODO())
		}
	}

	return nil
}

func findGroupForUserOption(user *apigateway.AuthenticatedUser, name string, link string) (*models.UserOption, error) {

	ok, coll := database.UserOptions()
	if !ok {
		return nil, errors.New("Not connected")
	}

	userMatch := bson.M{"userid": user.ViewID}
	nameOrLink := getNameOrLinkBson(name, link)

	filter := []bson.M{userMatch, nameOrLink}

	findOptions := options.Find()

	results, err := coll.Find(context.TODO(), bson.M{"$and": filter}, findOptions, &models.UserOption{})
	if err != nil {
		return nil, err
	}

	var bestOption *models.UserOption

	for i := range results {
		option := i.(*models.UserOption)
		if bestOption == nil {
			bestOption = option
		} else {
			if option.Link == link {
				bestOption = option
			}
		}
	}

	return bestOption, nil
}

func getNameOrLinkBson(name string, link string) bson.M {

	if link != "" {
		return bson.M{"$or": []bson.M{bson.M{"name": name}, bson.M{"link": link}}}
	}

	return bson.M{"name": name}
}

func (uo *UserOption) Remove(user *apigateway.AuthenticatedUser) error {

	if uo.db.UserID != user.ViewID {
		return errors.New("Cannot update another users options")
	}

	ok, coll := database.UserOptions()
	if !ok {
		return errors.New("Not connected")
	}

	return coll.DeleteByID(context.Background(), uo.db.ID)
}

func (gameSetup *GameSetUp) removeUserOptions(user *apigateway.AuthenticatedUser) error {

	gameSetUpID := bson.M{"gameSetupId": gameSetup.db.ID}
	idMatch := bson.M{"userid": user.ViewID}

	andBson := []bson.M{gameSetUpID, idMatch}

	ok, coll := database.UserOptions()
	if !ok {
		return errors.New("Not connected")
	}

	findOptions := options.Find()

	results, err := coll.Find(context.TODO(), andBson, findOptions, &models.UserOption{})
	if err != nil {
		return err
	}

	for i := range results {
		option := i.(*models.UserOption)
		coll.DeleteByID(context.Background(), option.ID)
	}

	return nil
}
