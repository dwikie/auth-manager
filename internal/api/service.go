package api

import (
	"github.com/dwikie/auth-manager/internal/handler"
	"github.com/dwikie/auth-manager/internal/service"
)

type Service struct {
	User *handler.UserHandler
}

func RegisterServices(repo *Repository) Service {
	userService := service.NewUserService(repo.User)

	return Service{
		User: handler.NewUserHandler(userService),
	}
}
