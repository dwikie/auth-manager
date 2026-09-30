package api

import "github.com/dwikie/auth-manager/internal/handler"

type Handler struct {
	User *handler.UserHandler
}

func InitializeHandler(service *Service) Handler {
	return Handler{
		User: service.User,
	}
}
