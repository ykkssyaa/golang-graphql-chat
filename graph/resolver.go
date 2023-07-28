package graph

import (
	"graphql_chat/package/model"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MessageChanals map[string]chan *model.Message
	Mutex          *sync.Mutex
}
