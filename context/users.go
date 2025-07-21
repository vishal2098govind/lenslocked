package context

import (
	"context"

	userM "github.com/vishal2098govind/lenslocked/models/user"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *userM.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *userM.User {
	user, ok := ctx.Value(userKey).(*userM.User)
	if !ok {
		return nil
	}
	return user
}
