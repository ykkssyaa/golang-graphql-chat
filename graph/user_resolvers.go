package graph

import (
	"context"
	"errors"
	"graphql_chat/package/common"
	"graphql_chat/package/model"
	"strconv"
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
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {

	customContext := common.GetContext(ctx)

	newUser := &model.UserDB{
		Name: input.Name,
	}

	err := customContext.Database.Create(newUser).Error

	return newUser.ToGraphQL(), err
}

// UserJoined is the resolver for the userJoined field.
func (r *subscriptionResolver) UserJoined(ctx context.Context, user string) (<-chan *model.Message, error) {

	userID, _ := strconv.ParseUint(user, 10, 64)
	err := common.GetContext(ctx).Database.First(&model.UserDB{}, uint(userID)).Error

	if err != nil {
		return nil, errors.New("Not enough user with id " + user)
	}

	chanal := make(chan *model.Message, 1)

	r.Mutex.Lock()
	r.MessageChanals[user] = chanal
	r.Mutex.Unlock()

	// Delete channel when done
	go func() {
		<-ctx.Done()
		r.Mutex.Lock()
		delete(r.MessageChanals, user)
		r.Mutex.Unlock()
	}()

	return chanal, nil
}
