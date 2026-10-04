package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/dwikie/auth-manager/internal/domain"
)

type ITokenService interface {
	GenerateAccessToken(claims *domain.TokenClaims, footer any) (string, error)
	GenerateRefreshToken(claims *domain.TokenClaims, footer any) (string, error)
}

type IUserService interface {
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	VerifyPassword(hashedPassword, password string) bool
}

type AuthenticateUserResponse struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	userService  IUserService
	tokenService ITokenService
}

func NewAuthService(userService IUserService, tokenService ITokenService) *AuthService {
	return &AuthService{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *AuthService) AuthenticateUser(ctx context.Context, username, password string) (*AuthenticateUserResponse, error) {
	user, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found") // User not found
	}

	if !s.userService.VerifyPassword(user.Password, password) {
		return nil, errors.New("invalid password") // Invalid password
	}

	encoder := base64.StdEncoding

	accessToken, err := s.tokenService.GenerateAccessToken(&domain.TokenClaims{
		JTI:        uuid.New().String(),
		Issuer:     "auth-manager",
		Audience:   "auth-manager",
		Subject:    encoder.EncodeToString(fmt.Append(nil, user.ID)),
		Expiration: time.Now().Add(time.Minute * 15),
		NotBefore:  time.Now(),
		IssuedAt:   time.Now(),
		Extensions: nil,
	}, nil)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(&domain.TokenClaims{
		JTI:        uuid.New().String(),
		Issuer:     "auth-manager",
		Audience:   "auth-manager",
		Subject:    encoder.EncodeToString(fmt.Append(nil, user.ID)),
		Expiration: time.Now().Add(time.Hour * 24),
		NotBefore:  time.Now(),
		IssuedAt:   time.Now(),
		Extensions: nil,
	}, nil)
	if err != nil {
		return nil, err
	}

	return &AuthenticateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil // Authentication successful
}
