package graph

import (
	"context"
	"fmt"
	"graphql_chat/package/common"
	"graphql_chat/package/model"
)

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	customContext := common.GetContext(ctx)

	var users []model.UserDB

	err := customContext.Database.Find(&users).Error

	if err != nil {
		return nil, err
	}

	res := make([]*model.User, len(users))

	for i, user := range users {
		res[i] = user.ToGraphQL()
	}

	return res, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input *model.NewUser) (*model.User, error) {

	customContext := common.GetContext(ctx)

	newUser := &model.UserDB{
		Name: input.Name,
	}

	err := customContext.Database.Create(newUser).Error

	return newUser.ToGraphQL(), err
}

// UserJoined is the resolver for the userJoined field.
func (r *subscriptionResolver) UserJoined(ctx context.Context, user string) (<-chan *model.User, error) {
	panic(fmt.Errorf("not implemented: UserJoined - userJoined"))
}
