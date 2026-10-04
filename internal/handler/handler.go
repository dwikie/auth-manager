package handler

import (
	"github.com/dwikie/auth-manager/internal/handler/auth"
	"github.com/dwikie/auth-manager/internal/handler/user"
	"github.com/dwikie/auth-manager/internal/service"
)

type Handler struct {
	UserHandler *user.UserHandler
	AuthHandler *auth.AuthHandler
}

func Initialize(service *service.Service) *Handler {
	userHandler := user.NewUserHandler(service.User)
	authHandler := auth.NewAuthHandler(service.Auth)

	return &Handler{
		UserHandler: userHandler,
		AuthHandler: authHandler,
	}
}
