package middleware

import (
	"context"
	"net/http"
	"strings"

	"lunikissShop/pkg/domain/model"
	"lunikissShop/pkg/domain/service"
)

type contextKey string

const (
	UserContextKey   contextKey = "user"
	UserIDContextKey contextKey = "user_id"
)

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				ctx := context.WithValue(r.Context(), UserContextKey, &model.User{
					Role: string(model.RoleAnonymous),
				})
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				tokenString = strings.TrimPrefix(authHeader, "bearer ")
			}

			user, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserContextKey, user)
			ctx = context.WithValue(ctx, UserIDContextKey, user.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*model.User)
	return user, ok
}
