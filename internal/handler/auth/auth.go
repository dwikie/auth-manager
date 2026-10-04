package auth

import (
	"github.com/dwikie/auth-manager/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type AuthenticateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthenticateUserResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *AuthHandler) AuthenticateUser(g *gin.Context) {
	var request AuthenticateUserRequest
	if err := g.ShouldBindJSON(&request); err != nil {
		g.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.AuthenticateUser(g.Request.Context(), request.Username, request.Password)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.SetCookie("refresh_token", token.RefreshToken, 3600, "/", "", true, true)

	g.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
		},
		"message": "User authenticated successfully",
		"errors":  nil,
	})
}
