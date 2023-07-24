package graph

import (
	"context"
	"errors"
	"fmt"
	"graphql_chat/package/common"
	"graphql_chat/package/model"
	"strconv"
)

// DeleteChat is the resolver for the deleteChat field.
func (r *mutationResolver) DeleteChat(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteChat - deleteChat"))
}

// Chats is the resolver for the chats field.
func (r *queryResolver) Chats(ctx context.Context, user *string) ([]*model.Chat, error) {
	customContext := common.GetContext(ctx)

	var chats []model.ChatDB
	var err error

	if user == nil {
		err = customContext.Database.Find(&chats).Error
	} else {
		userID, _ := strconv.ParseUint(*user, 10, 64)
		err = customContext.Database.Where("user1_id = ? OR user2_id = ?",
			userID, userID).Find(&chats).Error
	}

	if err != nil {
		return nil, err
	}

	res := make([]*model.Chat, len(chats))

	for i, chat := range chats {
		res[i] = chat.ToGraphQL()
	}

	return res, nil
}

// CreateChat is the resolver for the createChat field.
func (r *mutationResolver) CreateChat(ctx context.Context, input model.NewChat) (*model.Chat, error) {

	// TODO: Запрет создания чатов с одними и теми же пользователями

	if input.User1 == input.User2 {
		return nil, errors.New("chat with yourself")
	}

	customContext := common.GetContext(ctx)

	var user1, user2 model.UserDB

	id1, _ := strconv.ParseUint(input.User1, 10, 64)
	err := customContext.Database.First(&user1, uint(id1)).Error

	if err != nil {
		return nil, err
	}

	id2, _ := strconv.ParseUint(input.User2, 10, 64)
	err = customContext.Database.First(&user2, uint(id2)).Error

	if err != nil {
		return nil, err
	}

	newChat := &model.ChatDB{
		User1ID: uint(id1),
		User2ID: uint(id2),
	}

	err = customContext.Database.Omit("User1", "User2").Create(newChat).Error

	return newChat.ToGraphQL(), err
}
