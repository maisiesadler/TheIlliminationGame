package gamegroup

import (
	"context"
	"errors"

	"github.com/maisiesadler/theilliminationgame/apigateway"
	"github.com/maisiesadler/theilliminationgame/database"
)

func Create(ctx context.Context, user *apigateway.AuthenticatedUser, groupName string) (*Group, error) {
	ok, coll := database.Groups()
	if !ok {
		return nil, errors.New("Not connected")
	}

	members := []GroupMember{GroupMember{ID: user.ViewID, Nickname: user.Nickname}}
	dbgroup := &dbGroup{
		Members: members,
		Name:    groupName,
	}

	saved, err := coll.InsertOneAndFind(ctx, dbgroup, &dbGroup{})

	return loadGroup(saved.(*dbGroup)), err
}
