package context

import (
	"context"

	"github.com/jackytck/lenslocked/models"
)

const (
	userKey privateKey = "user"
)

type privateKey string

// WithUser sets the user in context.
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User retrieves user from context.
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
