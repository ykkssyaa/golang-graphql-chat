package graph

import (
	"context"
	"encoding/base64"
	"errors"
	"gorm.io/gorm/clause"
	"graphql_chat/package/common"
	"graphql_chat/package/model"
	"strconv"
	"strings"
	"time"
)

// PostMessage is the resolver for the postMessage field.
func (r *mutationResolver) PostMessage(ctx context.Context, input model.NewMessage) (*model.Message, error) {

	customContext := common.GetContext(ctx)

	senderID, _ := strconv.ParseUint(input.Sender, 10, 64)
	receiverID, _ := strconv.ParseUint(input.Receiver, 10, 64)

	var chat model.ChatDB

	err := customContext.Database.Where("user1_id = ? AND user2_id = ? OR user2_id = ? AND user1_id = ?",
		senderID, receiverID, senderID, receiverID).First(&chat).Error

	if err != nil {

		if err.Error() == "record not found" {

			createChatRes, err := r.CreateChat(ctx, model.NewChat{User1: input.Sender, User2: input.Receiver})

			if err != nil {
				return nil, err
			}

			chatId, _ := strconv.ParseUint(createChatRes.ID, 10, 64)
			chat.ID = uint(chatId)

		} else {
			return nil, err
		}
	}

	newMessage := model.MessageDB{
		Payload:    input.Payload,
		ChatID:     chat.ID,
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
	}

	if input.Time != nil {
		newMessage.Model.CreatedAt = *input.Time
	}

	// Сохранение в базу данных
	err = customContext.Database.Create(&newMessage).Error

	if err != nil {
		return nil, err
	}

	customContext.Database.Preload("Sender").Preload("Receiver").First(&newMessage)

	// Если пользователь подключен к серверу, то передаем ему в канал сообщение
	go func() {
		if userCanal, ok := r.MessageChanals[input.Receiver]; ok {
			r.Mutex.Lock()
			userCanal <- newMessage.ToGraphQL()
			r.Mutex.Unlock()
		}
	}()

	return newMessage.ToGraphQL(), nil
}

// DeleteMessage is the resolver for the deleteMessage field.
func (r *mutationResolver) DeleteMessage(ctx context.Context, id string) (*bool, error) {

	mesID, _ := strconv.ParseUint(id, 10, 64)
	var message model.MessageDB

	err := common.GetContext(ctx).Database.Preload("Receiver").Preload("Sender").
		Clauses(clause.Returning{}).Delete(&message, mesID).Error

	ok := true
	if err != nil {
		ok = false
		return &ok, err
	}

	if message.ID == 0 {
		ok = false
		return &ok, errors.New("message not exist")
	}

	// Если пользователь подключен к серверу, то передаем ему в канал удаленное сообщение
	go func() {

		mes := message.ToGraphQL()

		mes.Time = &time.Time{}

		if userCanal, ok := r.MessageChanals[strconv.FormatUint(uint64(message.ReceiverID), 10)]; ok {
			r.Mutex.Lock()
			userCanal <- mes
			r.Mutex.Unlock()
		}
	}()

	return &ok, err
}

// MessagesFromUser is the resolver for the messagesFromUser field.
func (r *queryResolver) MessagesFromUser(ctx context.Context, input model.MessagesFromUserInput, first *int, after *string) (*model.MessageConnection, error) {

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

	err := customContext.Database.Preload("Receiver").Preload("Sender").
		Where("sender_id = ? AND receiver_id = ?", senderID, receiverID).Find(&messagesFromUser).Error

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
