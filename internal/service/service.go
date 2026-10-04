package service

import (
	"github.com/dwikie/auth-manager/internal/config"
	"github.com/dwikie/auth-manager/internal/repo"
)

type Service struct {
	User  *UserService
	Auth  *AuthService
	Token *TokenService
}

func Initialize(repo *repo.Repository, secretProvider *config.Secrets) *Service {

	userService := NewUserService(repo.User)
	tokenService := NewTokenService(secretProvider)
	authService := NewAuthService(userService, tokenService)

	return &Service{
		User:  userService,
		Auth:  authService,
		Token: tokenService,
	}
}
