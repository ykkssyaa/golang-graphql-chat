package model

import "gorm.io/gorm"

type UserDB struct {
	gorm.Model
	name string `gorm:"not null;size:256"`
}

type MessageDB struct {
	gorm.Model
	payload    string `gorm:"not null;type:text"`
	chat       ChatDB // Many-to-one association
	senderID   uint
	receiverID uint
}

type ChatDB struct {
	gorm.Model
	user1 UserDB `gorm:"many2many:users_chats;"`
	user2 UserDB `gorm:"many2many:users_chats;"`
}
