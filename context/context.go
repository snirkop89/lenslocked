package context

import (
	"context"

	"github.com/snirkop89/lenslocked/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User)
	if !ok {
		// Protect against cases assertion will fail:
		// 1. If nothing was store in the first place
		// 2. Invalid type
		return nil
	}
	return user
}
