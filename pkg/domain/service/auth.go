package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"lunikissShop/pkg/domain/model"
)

type AuthService struct {
	userService          *UserService
	jwtSecret            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthService(userService *UserService, jwtSecret string) *AuthService {
	return &AuthService{
		userService:          userService,
		jwtSecret:            []byte(jwtSecret),
		accessTokenDuration:  24 * time.Hour,
		refreshTokenDuration: 7 * 24 * time.Hour,
	}
}

func (as *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	existingUser, err := as.userService.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser.ID != "" {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := as.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userInfo := &model.UserInfo{
		Email:    req.Email,
		Name:     req.Name,
		Surname:  req.Surname,
		Password: hashedPassword,
		Role:     string(model.RoleUser),
		Phone:    req.Phone,
	}

	if err := as.userService.AddUser(ctx, userInfo); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user, err := as.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get created user: %w", err)
	}

	accessToken, expiresAt, err := as.generateAccessToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthResponse{
		User:        user,
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func (as *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := as.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	hashedPassword, err := as.userService.GetUserPassword(ctx, user.ID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !as.checkPasswordHash(req.Password, hashedPassword) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, expiresAt, err := as.generateAccessToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthResponse{
		User:        user,
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func (as *AuthService) ValidateToken(tokenString string) (*model.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return as.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*model.TokenClaims); ok && token.Valid {
		return &model.User{
			ID:    claims.UserID,
			Email: claims.Email,
			Role:  string(claims.Role),
		}, nil
	}

	return nil, errors.New("invalid token")
}

func (as *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	user, err := as.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	fullUser, err := as.userService.GetUser(ctx, user.ID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, expiresAt, err := as.generateAccessToken(&fullUser)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &model.AuthResponse{
		User:        fullUser,
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func (as *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	currentHash, err := as.userService.GetUserPassword(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !as.checkPasswordHash(oldPassword, currentHash) {
		return errors.New("incorrect old password")
	}

	newHash, err := as.hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return as.userService.UpdatePassword(ctx, userID, newHash)
}

func (as *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (as *AuthService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (as *AuthService) generateAccessToken(user *model.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(as.accessTokenDuration)

	claims := &model.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   model.Role(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(as.jwtSecret)

	return tokenString, expiresAt, err
}

func (as *AuthService) generateRefreshToken(user *model.User) (string, error) {
	claims := &model.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   model.Role(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(as.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(as.jwtSecret)
}
