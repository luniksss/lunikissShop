package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
	Surname  string `json:"surname" validate:"required,min=2"`
	Phone    string `json:"phone" validate:"optional,min=11"`
}

type AuthResponse struct {
	User        User      `json:"user"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   Role   `json:"role"`
	jwt.RegisteredClaims
}
