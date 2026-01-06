package middleware

import (
	"context"
	"lunikissShop/pkg/domain/model"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*model.User)
	return user, ok
}
