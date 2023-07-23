package graph

import (
	"context"
	"fmt"
	"graphql_chat/package/common"
	"graphql_chat/package/model"
)

// DeleteChat is the resolver for the deleteChat field.
func (r *mutationResolver) DeleteChat(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteChat - deleteChat"))
}

// Chats is the resolver for the chats field.
func (r *queryResolver) Chats(ctx context.Context) ([]*model.Chat, error) {
	customContext := common.GetContext(ctx)

	var chats []model.ChatDB

	err := customContext.Database.Find(&chats).Error

	if err != nil {
		return nil, err
	}

	res := make([]*model.Chat, len(chats))

	for i, chat := range chats {
		res[i] = chat.ToGraphQL()
	}

	return res, nil
}
