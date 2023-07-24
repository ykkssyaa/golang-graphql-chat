package graph

import (
	"context"
	"encoding/base64"
	"fmt"
	"graphql_chat/package/common"
	"graphql_chat/package/model"
	"strconv"
	"strings"
)

// PostMessage is the resolver for the postMessage field.
func (r *mutationResolver) PostMessage(ctx context.Context, input model.NewMessage) (*model.Message, error) {

	customContext := common.GetContext(ctx)

	chatID, _ := strconv.ParseUint(input.Chat, 10, 64)
	senderID, _ := strconv.ParseUint(input.Sender, 10, 64)
	receiverID, _ := strconv.ParseUint(input.Receiver, 10, 64)

	newMessage := model.MessageDB{
		Payload:    input.Payload,
		ChatID:     uint(chatID),
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
	}

	err := customContext.Database.Create(&newMessage).Error

	if err != nil {
		return nil, err
	}

	return newMessage.ToGraphQL(), nil
}

// DeleteMessage is the resolver for the deleteMessage field.
func (r *mutationResolver) DeleteMessage(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteMessage - deleteMessage"))
}

// MessagePosted is the resolver for the messagePosted field.
func (r *subscriptionResolver) MessagePosted(ctx context.Context, message string) (<-chan *model.Message, error) {
	panic(fmt.Errorf("not implemented: MessagePosted - messagePosted"))
}

// MessagesFromUser is the resolver for the messagesFromUser field.
func (r *queryResolver) MessagesFromUser(ctx context.Context, input model.MessagesFromUserInput, first *int, after *string) (*model.MessageConnection, error) {

	// TODO: вложенность ответа

	customContext := common.GetContext(ctx)

	senderID, _ := strconv.ParseUint(input.Sender, 10, 64)
	receiverID, _ := strconv.ParseUint(input.Receiver, 10, 64)

	var from = 0

	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)

		if err != nil {
			return nil, err
		}

		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))

		if err != nil {
			return nil, err
		}

		from = i
	}

	var messagesFromUser []model.MessageDB

	err := customContext.Database.Where("sender_id = ? AND receiver_id = ?", senderID, receiverID).Find(&messagesFromUser).Error

	if err != nil {
		return nil, err
	}
	to := len(messagesFromUser)

	if first != nil {
		to = from + *first

		if to > len(messagesFromUser) {
			to = len(messagesFromUser)
		}
	}

	res := make([]*model.Message, len(messagesFromUser))

	for i, message := range messagesFromUser {
		res[i] = message.ToGraphQL()
	}

	return &model.MessageConnection{
		Messages: res,
		From:     from,
		To:       to,
	}, nil
}

// Edges is the resolver for the edges field.
func (r *messageConnectionResolver) Edges(ctx context.Context, obj *model.MessageConnection) ([]*model.MessageEdge, error) {
	edges := make([]*model.MessageEdge, obj.To-obj.From)

	for i := range edges {
		edges[i] = &model.MessageEdge{
			Node:   obj.Messages[obj.From+i],
			Cursor: model.EncodeCursor(obj.From + i),
		}
	}

	return edges, nil
}

// MessageConnection returns MessageConnectionResolver implementation.
func (r *Resolver) MessageConnection() MessageConnectionResolver {
	return &messageConnectionResolver{r}
}
