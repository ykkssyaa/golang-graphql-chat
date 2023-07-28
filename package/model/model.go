package model

import (
	"gorm.io/gorm"
	"strconv"
)

type UserDB struct {
	gorm.Model
	Name  string   `gorm:"not null;size:256"`
	chats []ChatDB `gorm:"many2many:users_chats;foreignKey:User1ID,User2ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *UserDB) ToGraphQL() *User {
	id := strconv.FormatUint(uint64(u.ID), 10)
	return &User{Name: u.Name, ID: id}
}

type MessageDB struct {
	gorm.Model
	Payload    string `gorm:"not null;type:text"`
	ChatID     uint
	SenderID   uint
	Sender     UserDB `gorm:"References:ID"`
	ReceiverID uint
	Receiver   UserDB `gorm:"References:ID"`
}

func (m *MessageDB) ToGraphQL() *Message {
	return &Message{
		ID:       strconv.FormatUint(uint64(m.ID), 10),
		Payload:  m.Payload,
		ChatID:   strconv.FormatUint(uint64(m.ChatID), 10),
		Receiver: m.Receiver.ToGraphQL(),
		Sender:   m.Sender.ToGraphQL(),
		Time:     &m.Model.CreatedAt,
	}
}

type ChatDB struct {
	gorm.Model
	Messages []MessageDB `gorm:"foreignKey:ChatID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User1    UserDB
	User2    UserDB
	User1ID  uint `gorm:"index:,unique,composite:idx_users_in_chat"`
	User2ID  uint `gorm:"index:,unique,composite:idx_users_in_chat"`
}

func (c *ChatDB) ToGraphQL() *Chat {
	return &Chat{
		ID:    strconv.FormatUint(uint64(c.ID), 10),
		User1: c.User1.ToGraphQL(),
		User2: c.User2.ToGraphQL(),
	}
}
