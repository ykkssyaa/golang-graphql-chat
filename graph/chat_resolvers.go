package graph

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"graphql_chat/package/common"
	"graphql_chat/package/model"
	"strconv"
	"time"
)

// DeleteChat is the resolver for the deleteChat field.
func (r *mutationResolver) DeleteChat(ctx context.Context, id string) (*bool, error) {

	chatID, _ := strconv.ParseUint(id, 10, 64)
	ok := true
	var chat model.ChatDB

	commonContext := common.GetContext(ctx)

	err := commonContext.Database.Transaction(func(tx *gorm.DB) error {
		err := commonContext.Database.Preload("User1").Preload("User2").
			Clauses(clause.Returning{}).Delete(&chat, chatID).Error

		if err != nil {
			return err
		}

		err = commonContext.Database.Where("chat_id = ?", chatID).Delete(&model.MessageDB{}).Error

		return err
	})

	if err != nil {
		ok = false
		return &ok, err
	}

	// Отправляем пользователям сигнал о том, что был удален чат
	go func() {

		// TODO: исправить отправку сообщения о том, что был удален чат

		mes := &model.Message{
			ChatID: id,
		}

		mes.Time = &time.Time{}

		user1ID := strconv.FormatUint(uint64(chat.User1ID), 10)
		user2ID := strconv.FormatUint(uint64(chat.User2ID), 10)

		if userCanal, ok := r.MessageChanals[user1ID]; ok {
			r.Mutex.Lock()
			userCanal <- mes
			r.Mutex.Unlock()
		}

		if userCanal, ok := r.MessageChanals[user2ID]; ok {
			r.Mutex.Lock()
			userCanal <- mes
			r.Mutex.Unlock()
		}
	}()

	return &ok, err

}

// Chats is the resolver for the chats field.
func (r *queryResolver) Chats(ctx context.Context, user *string) ([]*model.Chat, error) {
	customContext := common.GetContext(ctx)

	var chats []model.ChatDB
	var err error

	if user == nil {
		err = customContext.Database.Preload("User1").Preload("User2").Find(&chats).Error
	} else {
		userID, _ := strconv.ParseUint(*user, 10, 64)

		err = customContext.Database.Preload("User1").Preload("User2").
			Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&chats).Error

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
func (r *mutationResolver) CreateChat(ctx context.Context, input model.NewChat) (*model.ChatMutationResult, error) {

	if input.User1 == input.User2 {
		return nil, errors.New("chat with yourself")
	}

	customContext := common.GetContext(ctx)

	var user1, user2 model.UserDB

	id1, _ := strconv.ParseUint(input.User1, 10, 64)
	id2, _ := strconv.ParseUint(input.User2, 10, 64)

	// Предотвращаю создание чатов с одними и теми же пользователями по два раза
	if id1 > id2 {
		id1, id2 = id2, id1
	}

	err := customContext.Database.First(&user1, uint(id1)).Error

	if err != nil {
		return nil, err
	}

	err = customContext.Database.First(&user2, uint(id2)).Error

	if err != nil {
		return nil, err
	}

	newChat := &model.ChatDB{
		User1ID: uint(id1),
		User2ID: uint(id2),
	}

	err = customContext.Database.Omit("User1", "User2").Create(newChat).Error

	if err != nil {

		// TODO: Поймать ошибку SQL 23505
		if errors.Is(err, gorm.ErrForeignKeyViolated) {

			if err = customContext.Database.Model(&model.ChatDB{}).
				Where("user1_id = ? AND user2_id = ? OR user2_id = ? AND user1_id = ?",
					input.User1, input.User2, input.User2, input.User1).
				Update("deleted_at", nil).Error; err != nil {

				return nil, err
			}

		} else {
			return nil, err
		}
	}

	err = customContext.Database.Preload("User1").Preload("User2").First(&newChat, newChat.ID).Error

	return &model.ChatMutationResult{
		ID:    strconv.FormatUint(uint64(newChat.ID), 10),
		User1: strconv.FormatUint(uint64(newChat.User1ID), 10),
		User2: strconv.FormatUint(uint64(newChat.User2ID), 10),
	}, err
}
