package model

import (
	"gorm.io/gorm"
	"strconv"
)

type UserDB struct {
	gorm.Model
	Name string `gorm:"not null;size:256"`
}

func (u *UserDB) ToGraphQL() *User {
	id := strconv.FormatUint(uint64(u.ID), 10)
	return &User{Name: u.Name, ID: id}
}

type MessageDB struct {
	gorm.Model
	Payload    string `gorm:"not null;type:text"`
	Chat       ChatDB `gorm:"foreignKey:ID"` // Many-to-one association
	SenderID   uint
	ReceiverID uint
}

type ChatDB struct {
	gorm.Model
	User1 UserDB `gorm:"many2many:users_chats;"`
	User2 UserDB `gorm:"many2many:users_chats;"`
}
