package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.Engine, handler *Handler) {
	api := g.Group("/api")

	userRoutes := api.Group("/users")
	userRoutes.POST("/", handler.User.CreateUser)
}
