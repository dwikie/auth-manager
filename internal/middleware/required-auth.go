package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RequiredAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || !strings.HasPrefix(parts[1], "v2.local.") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			return
		}

		at := parts[1]

		claims, _, errs := m.tokenService.ParseAccessToken(at)
		fmt.Println("Access Token Claims: ", claims)
		if errs != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": errs.Error(),
			})
			return
		}

		useridStr, err := base64.StdEncoding.DecodeString(claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token subject"})
			return
		}

		userid, err := strconv.Atoi(string(useridStr))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token subject"})
			return
		}

		c.Set("userid", userid)
		c.Next()
	}
}
