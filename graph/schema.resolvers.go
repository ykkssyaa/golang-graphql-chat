package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"
	"fmt"
	"graphql_chat/graph/model"
)

// Edges is the resolver for the edges field.
func (r *messageConnectionResolver) Edges(ctx context.Context, obj *model.MessageConnection) ([]*model.MessageEdge, error) {
	panic(fmt.Errorf("not implemented: Edges - edges"))
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input *model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}

// PostMessage is the resolver for the postMessage field.
func (r *mutationResolver) PostMessage(ctx context.Context, input *model.NewMessage) (*model.Message, error) {
	panic(fmt.Errorf("not implemented: PostMessage - postMessage"))
}

// DeleteChat is the resolver for the deleteChat field.
func (r *mutationResolver) DeleteChat(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteChat - deleteChat"))
}

// DeleteMessage is the resolver for the deleteMessage field.
func (r *mutationResolver) DeleteMessage(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteMessage - deleteMessage"))
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented: Users - users"))
}

// Chats is the resolver for the chats field.
func (r *queryResolver) Chats(ctx context.Context) ([]*model.Chat, error) {
	panic(fmt.Errorf("not implemented: Chats - chats"))
}

// MessagesFromUser is the resolver for the messagesFromUser field.
func (r *queryResolver) MessagesFromUser(ctx context.Context, input model.MessagesFromUserInput, first *int, after *string) (*model.MessageConnection, error) {
	panic(fmt.Errorf("not implemented: MessagesFromUser - messagesFromUser"))
}

// MessagePosted is the resolver for the messagePosted field.
func (r *subscriptionResolver) MessagePosted(ctx context.Context, message string) (<-chan *model.Message, error) {
	panic(fmt.Errorf("not implemented: MessagePosted - messagePosted"))
}

// UserJoined is the resolver for the userJoined field.
func (r *subscriptionResolver) UserJoined(ctx context.Context, user string) (<-chan *model.User, error) {
	panic(fmt.Errorf("not implemented: UserJoined - userJoined"))
}

// MessageConnection returns MessageConnectionResolver implementation.
func (r *Resolver) MessageConnection() MessageConnectionResolver {
	return &messageConnectionResolver{r}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type messageConnectionResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
