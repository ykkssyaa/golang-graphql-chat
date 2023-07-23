package graph

import (
	"context"
	"fmt"
	"graphql_chat/package/model"
)

// Edges is the resolver for the edges field.
func (r *messageConnectionResolver) Edges(ctx context.Context, obj *model.MessageConnection) ([]*model.MessageEdge, error) {
	panic(fmt.Errorf("not implemented: Edges - edges"))
}

// PostMessage is the resolver for the postMessage field.
func (r *mutationResolver) PostMessage(ctx context.Context, input *model.NewMessage) (*model.Message, error) {
	panic(fmt.Errorf("not implemented: PostMessage - postMessage"))
}

// DeleteMessage is the resolver for the deleteMessage field.
func (r *mutationResolver) DeleteMessage(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteMessage - deleteMessage"))
}

// MessagesFromUser is the resolver for the messagesFromUser field.
func (r *queryResolver) MessagesFromUser(ctx context.Context, input model.MessagesFromUserInput, first *int, after *string) (*model.MessageConnection, error) {
	panic(fmt.Errorf("not implemented: MessagesFromUser - messagesFromUser"))
}

// MessagePosted is the resolver for the messagePosted field.
func (r *subscriptionResolver) MessagePosted(ctx context.Context, message string) (<-chan *model.Message, error) {
	panic(fmt.Errorf("not implemented: MessagePosted - messagePosted"))
}

// MessageConnection returns MessageConnectionResolver implementation.
func (r *Resolver) MessageConnection() MessageConnectionResolver {
	return &messageConnectionResolver{r}
}
