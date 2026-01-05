package service

import (
	"context"
	"errors"

	"lunikissShop/pkg/domain/model"
	"lunikissShop/pkg/infrastructure/mysql/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo}
}

func (us UserService) ListAllUsers(ctx context.Context) ([]model.User, error) {
	return us.userRepo.ListAllUsers(ctx)
}

func (us UserService) GetUser(ctx context.Context, userID string) (model.User, error) {
	return us.userRepo.GetUserByID(ctx, userID)
}

func (us UserService) GetUserByEmail(ctx context.Context, userEmail string) (model.User, error) {
	return us.userRepo.GetUserByEmail(ctx, userEmail)
}

func (us UserService) AddUser(ctx context.Context, userInfo *model.UserInfo) error {
	_, err := us.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		return errors.New("user already exists")
	}

	return us.userRepo.AddUser(ctx, userInfo)
}

func (us UserService) UpdateUser(ctx context.Context, newUserInfo *model.UserInfo) error {
	_, err := us.GetUser(ctx, newUserInfo.ID)
	if err != nil {
		return errors.New("user not exists")
	}

	return us.userRepo.UpdateUser(ctx, newUserInfo)
}

func (us UserService) UpdateUserRole(ctx context.Context, userID string, newUserRole string) error {
	_, err := us.GetUser(ctx, userID)
	if err != nil {
		return errors.New("user not exists")
	}

	return us.userRepo.UpdateUserRole(ctx, userID, newUserRole)
}

func (us UserService) DeleteUser(ctx context.Context, userID string) error {
	_, err := us.GetUser(ctx, userID)
	if err != nil {
		return errors.New("user not exists")
	}

	return us.userRepo.DeleteUser(ctx, userID)
}
