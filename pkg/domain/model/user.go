package model

import "context"

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Phone   string `json:"phone,omitempty"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Phone    string `json:"phone,omitempty"`
}

type UserRepository interface {
	ListAllUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, userID string) (User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (User, error)
	AddUser(ctx context.Context, userInfo *UserInfo) error
	UpdateUser(ctx context.Context, newUserInfo *UserInfo) error
	UpdateUserRole(ctx context.Context, userID string, newUserRole string) error
	DeleteUser(ctx context.Context, userID string) error
}
