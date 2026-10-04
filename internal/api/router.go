package api

import (
	"github.com/dwikie/auth-manager/internal/handler"
	"github.com/dwikie/auth-manager/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.Engine, handler *handler.Handler, middleware *middleware.Middleware) {
	api := g.Group("/api")

	userRoutes := api.Group("/users")
	userRoutes.POST("/", handler.UserHandler.CreateUser)
	userRoutes.GET("/profile", middleware.RequiredAuthentication(), handler.UserHandler.GetCurrentUserProfile)

	authRoutes := api.Group("/auth")
	authRoutes.POST("/login", handler.AuthHandler.AuthenticateUser)
}
