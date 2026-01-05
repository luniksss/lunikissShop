package service

import (
	"context"
	"lunikissShop/pkg/domain/model"
)

type UserService interface {
	ListAllUsers(ctx context.Context) ([]model.User, error)
	GetUser(ctx context.Context, userID string) (model.User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (model.User, error)
	AddUser(ctx context.Context, userInfo *model.UserInfo) error
	UpdateUser(ctx context.Context, newUserInfo *model.UserInfo) error
	UpdateUserRole(ctx context.Context, userID string, newUserRole string) error
	DeleteUser(ctx context.Context, userID string) error
}
